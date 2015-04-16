package handles

import (
	"github.com/SaviorPhoenix/autobd/config"
	"github.com/SaviorPhoenix/autobd/dir"
	"golang.org/x/exp/inotify"
	"log"
)

type EventHandle func(*dir.Dir, *inotify.Event, string) error

func ParseMask(mask uint32) []string {
	var ret []string
	var flags = map[uint32]string{
		inotify.IN_CLOSE_WRITE: "changed",
		inotify.IN_CREATE:      "created",
		inotify.IN_DELETE:      "deleted",
		inotify.IN_MOVE:        "moved",
		inotify.IN_MOVED_FROM:  "moved_from",
		inotify.IN_MOVED_TO:    "moved_to",
		inotify.IN_OPEN:        "opened",
		inotify.IN_CLOSE:       "closed",
		inotify.IN_MOVE_SELF:   "watch_moved",
		inotify.IN_UNMOUNT:     "unmount",
	}
	for flag, _ := range flags {
		if flag&mask != 0 {
			ret = append(ret, flags[flag])
		}
	}
	return ret
}

func HandleEvent(handles map[string][]EventHandle, dir *dir.Dir, event *inotify.Event) {
	eventNames := ParseMask(event.Mask)

	//iterate through the names..
	for _, name := range eventNames {
		//..get the handle array
		arr := handles[name]
		//and finally iterate through the handle array calling each handle
		for _, handle := range arr {
			if handle != nil {
				handle(dir, event, name)
			} else {
				log.Println("No handles for event:")
			}
		}
	}
}

//Returns a map of arrays of event handles indexed by names instead of uint32 flags
func GetEventHandles(conf *config.Config) map[string][]EventHandle {
	var handles = map[string]EventHandle{
		"log":            LogEvent,
		"panic":          PanicEvent,
		"update_servers": UpdateServersEvent,
		"hash_files":     HashEvent,
		"nop":            NopEvent,
	}
	var ret map[string][]EventHandle

	ret = make(map[string][]EventHandle)
	for Name, arr := range conf.EventActions {
		for _, action := range arr {
			ret[Name] = append(ret[Name], handles[action])
		}
	}
	return ret
}

func HandleDirCreate(eventDir *dir.Dir, event *inotify.Event, name string) error {
	log.Printf("Adding child of %s: '%s' to tree", eventDir.GetPath(), name)
	return nil
}

func HandleDirRm(eventDir *dir.Dir, event *inotify.Event, name string) error {
	log.Printf("Removing child of %s: '%s' from tree", eventDir.GetPath(), name)
	return nil
}

func LogEvent(eventDir *dir.Dir, event *inotify.Event, name string) error {
	log.Println(eventDir.GetName() + ": " + name)
	return nil
}

func PanicEvent(eventDir *dir.Dir, event *inotify.Event, name string) error {
	log.Printf("PANIC EVENT: Event Directory: %s\nInotify Event: %s\n", eventDir, event)
	panic("")
}

func UpdateServersEvent(eventDir *dir.Dir, event *inotify.Event, name string) error {
	log.Println("UpdateServersEvent(): Not yet implemented")
	return nil
}

func HashEvent(eventDir *dir.Dir, event *inotify.Event, name string) error {
	log.Println("HashEvent(): Not yet implemented")
	return nil
}

func NopEvent(eventDir *dir.Dir, event *inotify.Event, name string) error {
	return nil
}
