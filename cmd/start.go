package main

import (
	"log"
	"redis_server"
)

func main() {
	server := redis_server.NewServer(redis_server.ServerConfig{Port: 6379})

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
