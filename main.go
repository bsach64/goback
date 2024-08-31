package main

import (
	"time"

	"github.com/bsach64/goback/client"
	"github.com/bsach64/goback/server"
	"log"
)

func main() {
	// f, err := utils.ChunkFile("example.txt")

	go server.Listen()
	time.Sleep(2 * time.Second)
	c, err := client.ConnectToServer("demo", "password", "127.0.0.1:2022")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer c.Close()
	client.Upload(c, "example.txt")
}

// Push this where
