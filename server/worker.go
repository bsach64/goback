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
		log.Fatalf("Cannot get the working directory of %v", err)
	}
	rsaPath := wd + "/private/id_rsa"

	err = utils.WatchDirectory("./.data/")
	if err != nil {
		log.Fatalf("Error while creating watcher %v", err)
	}

	sftpServer := New(w.Ip, rsaPath, w.Port)
	w.sftpServer = &sftpServer

	go func() {
		err := w.sftpServer.Listen()
		if err != nil {
			log.Fatalf("Worker SFTP server failed to listen: %v", err)
		}
	}()
}
