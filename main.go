package main

import (
	"log"
)

func main() {
	//router := gin.Default()
	//
	//router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	//

	//
	//router.Run(":8080")

	storage, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIService(":4500", storage)
	server.Run()
}
