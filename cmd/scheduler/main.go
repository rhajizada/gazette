package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/database"
)

var Version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Gazette scheduler %s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.LoadServer()
	if err != nil {
		log.Panicf("error loading config: %v", err)
	}
	conn := database.CreateRedisClient(&cfg.Redis)
	scheduler := asynq.NewScheduler(conn, nil)

	if err := scheduler.Run(); err != nil {
		log.Panicf("could not run scheduler: %v", err)
	}
}
