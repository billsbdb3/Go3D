package main

import (
	"3d-library/internal/database"
	"3d-library/internal/jobs"
	"log"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("✓ Worker connected to database")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 10},
	)

	mux := jobs.NewServer(db)

	log.Println("✓ Worker started, processing jobs...")
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
