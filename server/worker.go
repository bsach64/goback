package server

import (
	"log"
)

// Worker is just a usual SFTP server that handles the file request
// Master assigns client with a worker server
type Worker struct {
	id         int
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	SftpUser   string `json:"sftpUser"`
	SftpPass   string `json:"sftpPass"`
	sftpServer *SFTPServer
	master     *Server
}

func (w *Worker) StartSFTPServer() {
	sftpServer := New(w.Ip, "private/id_rsa", w.Port)
	w.sftpServer = &sftpServer

	go func() {
		err := w.sftpServer.Listen()
		if err != nil {
			log.Fatalf("Worker SFTP server failed to listen: %v", err)
		}
	}()
}
