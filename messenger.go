package main

import (
	"os"

	"github.com/hennersz/gomessenger/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	arguments := os.Args

	if len(arguments) == 1 {
		log.Fatal("Please provide a port number")
	}

	server := server.NewMessagingServer("127.0.0.1:" + arguments[1])
	server.Run()
}
