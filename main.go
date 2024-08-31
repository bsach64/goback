package main

import (
	// "fmt"
	"time"

	"github.com/bsach64/goback/client"
	"github.com/bsach64/goback/server"
	// "github.com/bsach64/goback/utils"
)

func main() {
	// f, err := utils.ChunkFile("example.txt")

	go server.Listen()
	time.Sleep(2 * time.Second)
	client.Upload("example.txt")

}

// Push this where
