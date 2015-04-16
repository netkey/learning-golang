package main

import (
	"log"

	"github.com/go-fsnotify/fsnotify"
)

func myNewWatcher(watchDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main() {

	watchDir := "/Users/qinshen/git/project/tsingcloud/bin"

	myNewWatcher(watchDir)

	for {
	}
}
