package utils

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

func WatchDirectory(directory string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error while creating directory watcher: %v", err)
	}

	log.Info("Starting to watch directory:", directory)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Error("Watcher events channel closed")
					return
				}
				log.Info("Event: ", event)
				if event.Has(fsnotify.Write) {
					log.Info("File modified: ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Error("Watcher errors channel closed")
					return
				}
				log.Errorf("Watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		watcher.Close() // Clean up if watcher.Add fails
		return nil, fmt.Errorf("failed to add directory to watcher: %v", err)
	}

	return watcher, nil
}
