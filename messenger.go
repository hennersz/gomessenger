package main

import (
	"log"
	"os"

	"github.com/hennersz/gomessenger/server"
)

func main() {
	arguments := os.Args

	if len(arguments) == 1 {
		log.Fatal("Please provide a port number")
	}

	portNum := arguments[1]
	server.StartServer(portNum)
}
