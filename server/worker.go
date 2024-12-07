package server

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bsach64/goback/utils"
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
	configPath := filepath.Join(wd, "Directory.json")

	config := readConf(configPath)

	watchDir := filepath.Join(wd, config.Directory)

	if config.Directory != "" {
		watcher, err := utils.WatchDirectory(watchDir)
		if err != nil {
			log.Errorf("Error while creating watcher: %v", err)
			defer watcher.Close() // Ensure the watcher is cleaned up on server shutdown
		}
	} else {
		log.Info("No directory is being watched you can add using", "command", "Add Directory to Sync")
	}

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

type Config struct {
	Directory string `json:"dir"`
}

func readConf(path string) Config {
	config := Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.Directory = "./.data"
		file, err := os.Create(path)
		if err != nil {
			log.Fatalf("Failed to create config file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(config); err != nil {
			log.Fatalf("Failed to write default config to file: %v", err)
		}
		log.Printf("Config file created with default directory: %s", config.Directory)
	} else {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("Failed to open config file: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			log.Fatalf("Failed to parse config file: %v", err)
		}
	}

	return config
}
