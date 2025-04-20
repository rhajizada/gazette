package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/tasks"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/database"
	"github.com/rhajizada/gazette/internal/repository"
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

	cfg, err := config.LoadWorker()
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
	c := *asynq.NewClient(conn)
	err = c.Ping()
	if err != nil {
		log.Panicf("failed to connect to Redis: %v", err)
	}
	if err != nil {
		log.Panicf("failed to connect to Redis: %v", err)
	}

	srv := asynq.NewServer(conn, asynq.Config{})

	h := tasks.NewHandler(rq)
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeFeedSync, h.HandleFeedSync)

	if err := srv.Run(mux); err != nil {
		log.Panicf("could not run server: %v", err)
	}
}
