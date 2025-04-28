package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/database"
	"github.com/rhajizada/gazette/internal/handler"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/oauth"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/router"
	"github.com/rhajizada/gazette/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

var Version = "dev"

// @title Gazette API
// @description Swagger API documentation for Gazette.
// @version 0.1.0

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Gazette Server version %s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.LoadServer()
	if err != nil {
		log.Panicf("error loading config: %v", err)
	}

	pool, err := database.CreatePool(&cfg.Database)
	if err != nil {
		log.Panicf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	migrationsDir := "data/sql/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Panicf("migrations directory does not exist: %s", migrationsDir)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Panicf("failed to set goose dialect: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Panicf("failed to apply migrations: %v", err)
	}

	rq := repository.New(pool)
	conn := database.CreateRedisClient(&cfg.Redis)
	client := *asynq.NewClient(conn)
	err = client.Ping()
	if err != nil {
		log.Panicf("failed to connect to Redis: %v", err)
	}

	verifier, err := oauth.GetVerifier(&cfg.OAuth)
	if err != nil {
		log.Panicf("failed to initialize auth provider: %v", err)
	}

	oauthCfg, err := oauth.GetConfig(&cfg.OAuth)
	if err != nil {
		log.Panicf("failed to initialize auth provider: %v", err)
	}

	// Create handler
	service := service.New(rq, &client)
	handler := handler.New(service, []byte(cfg.SecretKey), verifier, oauthCfg)

	mux := http.NewServeMux()
	feedsAPIRoutes := router.RegisterFeedsAPI(handler)
	itemsAPIRoutes := router.RegisterItemsAPI(handler)
	collectionsAPIRoutes := router.RegisterCollectionsAPI(handler)
	oauthRoutes := router.RegisterOAuthRoutes(handler)

	loggingMiddleware := middleware.Logging()
	authMiddleware := middleware.APIAuthMiddleware([]byte(cfg.SecretKey))

	mux.Handle("/api/feeds/", http.StripPrefix("/api", authMiddleware(feedsAPIRoutes)))
	mux.Handle("/api/collections/", http.StripPrefix("/api", authMiddleware(collectionsAPIRoutes)))
	mux.Handle("/api/items/", http.StripPrefix("/api", authMiddleware(itemsAPIRoutes)))
	mux.Handle("/api/docs/", httpSwagger.WrapHandler)
	mux.Handle("/oauth/", http.StripPrefix("/oauth", oauthRoutes))
	mux.Handle("/", http.HandlerFunc(handler.WebHandler))

	log.Printf("server is running on port %v\n", cfg.Port)
	addr := fmt.Sprintf(":%v", cfg.Port)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Panicf("could not start server: %s\n", err.Error())
	}
}
