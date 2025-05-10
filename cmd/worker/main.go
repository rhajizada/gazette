package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/workers"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/database"
	"github.com/rhajizada/gazette/internal/repository"
)

var Version = "dev"

func main() {
	nullFile, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		log.Println("error opening /dev/null:", err)
		return
	}
	defer nullFile.Close()

	os.Stdout = nullFile
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Gazette Worker version %s\n", Version)
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
	client := *asynq.NewClient(conn)
	err = client.Ping()
	if err != nil {
		log.Panicf("failed to connect to Redis: %v", err)
	}
	if err != nil {
		log.Panicf("failed to connect to Redis: %v", err)
	}

	serverConfig := workers.GetConfig(&cfg.Queues)

	ollamaClient, err := workers.GetOllamaClient(&cfg.Ollama)
	if err != nil {
		log.Panicf("failed to initialize Ollama client: %v", err)
	}
	err = workers.InitModels(ollamaClient, &cfg.Ollama)
	if err != nil {
		log.Panicf("failed to initialize models: %v", err)
	}

	server := asynq.NewServer(conn, *serverConfig)

	handler := workers.NewHandler(rq, &client, &cfg.Ollama)
	mux := asynq.NewServeMux()
	mux.HandleFunc(workers.TypeSyncData, handler.HandleDataSync)
	mux.HandleFunc(workers.TypeSyncFeed, handler.HandleFeedSync)
	mux.HandleFunc(workers.TypeEmbedItem, handler.HandleEmbedItem)
	mux.HandleFunc(workers.TypeEmbedUser, handler.HandleEmbedUser)

	if err := server.Run(mux); err != nil {
		log.Panicf("could not run worker: %v", err)
	}
}
