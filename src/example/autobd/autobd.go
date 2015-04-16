package main

import (
	"fmt"
	"github.com/SaviorPhoenix/autobd/config"
	"github.com/SaviorPhoenix/autobd/dir"
	"github.com/SaviorPhoenix/autobd/handles"
	"golang.org/x/exp/inotify"
	"log"
	"os"
	"path/filepath"
)

const Logo = `
 █████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ██████╗
██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██╔══██╗
███████║██║   ██║   ██║   ██║   ██║██████╔╝██║  ██║
██╔══██║██║   ██║   ██║   ██║   ██║██╔══██╗██║  ██║
██║  ██║╚██████╔╝   ██║   ╚██████╔╝██████╔╝██████╔╝
╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝
Backing you up since whenever..
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<Path to watch dir>")
		return
	}

	fmt.Println(Logo)
	conf := config.ParseConfig("config.toml")

	eventHandles := handles.GetEventHandles(conf)
	watchFlags := conf.GetFlags()

	if eventHandles == nil || watchFlags == 0 {
		panic("Failed to get sane config")
	}

	rootPath := filepath.Clean(os.Args[1])
	rootName := filepath.Base(rootPath)

	root := dir.NewDir(rootPath, rootName, watchFlags, conf.Options.Recursive)

	if err := root.StartWatch(); err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-root.Watcher.Event:
			eventDir := root.GetChildDir(filepath.Clean(event.Name))

			if event.Mask&inotify.IN_ISDIR != 0 {
				parent, child := filepath.Split(filepath.Clean(event.Name))
				if eventDir == nil && event.Mask&inotify.IN_CREATE != 0 {
					eventDir := root.GetChildDir(filepath.Clean(parent))
					handles.HandleDirCreate(eventDir, event, child)
					continue
				} else if event.Mask&inotify.IN_DELETE != 0 {
					eventDir := root.GetChildDir(filepath.Clean(parent))
					handles.HandleDirRm(eventDir, event, child)
					continue
				}
			}
			handles.HandleEvent(eventHandles, eventDir, event)
		case err := <-root.Watcher.Error:
			log.Println(err)
		}
	}
}
