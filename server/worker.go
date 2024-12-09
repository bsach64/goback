package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/bsach64/goback/utils"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
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
	fileWrites := make(map[string]time.Time)
	delay := 10 * time.Millisecond

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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Error("Watcher events channel closed")
					return
				}
				if event.Has(fsnotify.Create | fsnotify.Write) {
					fileWrites[event.Name] = time.Now()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Error("Watcher errors channel closed")
					return
				}
				log.Errorf("Watcher error: %v", err)
			}

			for filename, t := range fileWrites {
				if time.Since(t) > delay {
					time.Sleep(delay)
				}
				err := reconstruct(wd, filename)
				if err != nil {
					log.Error("Could not reconstruct file", filename, "err", err)
					continue
				}
				log.Info("Reconstructed file", filename)
				delete(fileWrites, filename)
			}
		}
	}()

	err = os.MkdirAll(filepath.Join(wd, ".data", "snapshots"), 0755)
	if err != nil {
		watcher.Close()
		log.Fatalf("failed to create directory for watcher: %v", err)
	}

	err = watcher.Add(filepath.Join(wd, ".data", "snapshots"))
	if err != nil {
		watcher.Close() // Clean up if watcher.Add fails
		log.Fatalf("failed to add directory to watcher: %v", err)
	}

	// Keep the worker process alive
	select {}
}

func reconstruct(wd, filename string) error {
	dat, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var snapshot utils.Snapshot
	err = json.Unmarshal(dat, &snapshot)
	if err != nil {
		return err
	}

	byteData, err := utils.Reconstruct(snapshot)
	if err != nil {
		return err
	}

	log.Info("BYTES", string(byteData))

	err = os.MkdirAll(filepath.Join(wd, "files"), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(wd, snapshot.Filename))
	if err != nil {
		return err
	}

	_, err = file.Write(byteData)
	if err != nil {
		return err
	}
	log.Info("Recreated: ", "file", filename)
	return nil
}
