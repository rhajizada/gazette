package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/database"
	"github.com/rhajizada/gazette/internal/tasks"
)

var Version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Gazette Scheduler version %s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.LoadScheduler()
	if err != nil {
		log.Panicf("error loading config: %v", err)
	}
	conn := database.CreateRedisClient(&cfg.Redis)
	scheduler := asynq.NewScheduler(conn, &asynq.SchedulerOpts{
		HeartbeatInterval: cfg.HeartbeatInterval,
		Location:          cfg.Location,
	})

	err = scheduler.Ping()
	if err != nil {
		log.Panicf("failed to initialze scheduler: %v", err)
	}

	dataSyncTask, _ := tasks.NewDataSyncTask()

	id, err := scheduler.Register("@every 30m", dataSyncTask)
	if err != nil {
		log.Panicf("failed scheduling data sync task: %v", err)
	}
	log.Printf("scheduled data sync task %s", id)

	if err := scheduler.Run(); err != nil {
		log.Panicf("could not run scheduler: %v", err)
	}
}
