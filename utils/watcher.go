package utils

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func WatchDirectory(directory string) error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return fmt.Errorf("Error while file watcher")
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
