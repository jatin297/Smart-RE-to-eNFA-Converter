package main

import (
	redis2 "github.com/jatin297/retoenfa/redis"
	"log"
)

func main() {
	storage, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	client, err := redis2.NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIService(":4500", storage, client)
	server.Run()
}
