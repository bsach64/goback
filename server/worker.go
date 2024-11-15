package server

import (
	"log"

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

	err := utils.WatchDirectory("./.data/")
	if err != nil {
		log.Fatalf("Error while creating watcher %v", err)
	}
	sftpServer := New(w.Ip, "/app/private/id_rsa", w.Port)
	w.sftpServer = &sftpServer

	go func() {
		err := w.sftpServer.Listen()
		if err != nil {
			log.Fatalf("Worker SFTP server failed to listen: %v", err)
		}
	}()
}
