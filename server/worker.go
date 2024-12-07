package server

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
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
		log.Errorf("Cannot get the working directory: %v", err)
	}

	rsaPath := filepath.Join(wd, "private", "id_rsa")
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
