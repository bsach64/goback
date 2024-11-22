package server

import (
	"log"
	"os"

	"github.com/bsach64/goback/utils"
)

// Worker is just a usual SFTP server that handles the file request
// Master assigns client with a worker server
type Worker struct {
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	SftpUser   string `json:"sftpUser"`
	SftpPass   string `json:"sftpPass"`
	sftpServer *SFTPServer
}

func (w *Worker) StartSFTPServer() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get the working directory: %v", err)
	}

	rsaPath := wd + "/private/id_rsa"

	// Start directory watcher
	watcher, err := utils.WatchDirectory(wd + "/.data")
	if err != nil {
		log.Fatalf("Error while creating watcher: %v", err)
	}
	defer watcher.Close() // Ensure the watcher is cleaned up on server shutdown

	sftpServer := New(w.Ip, rsaPath, w.Port)
	w.sftpServer = &sftpServer

	go func() {
		err := w.sftpServer.Listen()
		if err != nil {
			log.Fatalf("Worker SFTP server failed to listen: %v", err)
		}
	}()

	// Keep the worker process alive
	select {}
}
