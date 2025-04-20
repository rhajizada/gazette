package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/database"
	"github.com/rhajizada/gazette/internal/handler"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/router"
)

var Version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	// If the version flag is provided, print version info and exit.
	if *versionFlag {
		fmt.Printf("Gazette %s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.Load()
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

	_ = repository.New(pool)

	rq := repository.New(pool)

	// Create handler
	h := handler.New(rq, []byte(cfg.SecretKey))

	r := http.NewServeMux()
	feeds := router.RegisterFeedRoutes(h)
	r.Handle("/api/feeds", http.StripPrefix("/api", feeds))
	lm := middleware.Logging(r)

	// Start the server
	log.Printf("server is running on port %v\n", cfg.Port)
	addr := fmt.Sprintf(":%v", cfg.Port)
	if err := http.ListenAndServe(addr, lm); err != nil {
		log.Panicf("could not start server: %s\n", err.Error())
	}
}
