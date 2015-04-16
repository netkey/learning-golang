package config

import (
	//	"github.com/SaviorPhoenix/autobd/handles"
	"github.com/BurntSushi/toml"
	"golang.org/x/exp/inotify"
	"log"
)

type Config struct {
	Options struct {
		Recursive      bool     `toml:"recursive"`
		InheritFlags   bool     `toml:"children_inherit_flags"`
		HashFile       bool     `toml:"hash_files"`
		UpdateInterval int      `toml:"server_update_interval"`
		Servers        []string `toml:"servers"`
		WatchFlags     []string `toml:"watch_flags"`
	} `toml:"Options"`
	EventActions map[string][]string `toml:"event_actions"`
}

func NewConfig() *Config {
	conf := &Config{}

	conf.Options.Recursive = true
	conf.Options.InheritFlags = true
	conf.Options.HashFile = true
	conf.Options.UpdateInterval = 5
	conf.Options.Servers = []string{""}
	conf.Options.WatchFlags = []string{"changed", "created", "deleted", "moved", "moved_from", "moved_to",
		"opened", "closed", "watch_moved", "unmount"}

	conf.EventActions["changed"] = []string{"update_servers", "log"}
	conf.EventActions["created"] = []string{"update_servers", "log"}
	conf.EventActions["deleted"] = []string{"update_servers", "log"}
	conf.EventActions["moved"] = []string{"update_servers", "log"}
	conf.EventActions["moved_from"] = []string{"update_servers", "log"}
	conf.EventActions["moved_to"] = []string{"update_servers", "log"}
	conf.EventActions["opened"] = []string{"log"}
	conf.EventActions["closed"] = []string{"log"}
	conf.EventActions["watch_moved"] = []string{"log"}
	conf.EventActions["unmount"] = []string{"log"}

	return conf
}

func ParseConfig(filename string) *Config {
	var conf *Config
	if _, err := toml.DecodeFile(filename, &conf); err != nil {
		log.Println(err)
		log.Println("Could not read or parse config, using defaults instead")
		return NewConfig()
	} else {
		log.Println("Got config from", filename)
	}
	return conf
}

func (conf *Config) GetFlags() uint32 {
	var flags = map[string]uint32{
		"created":     inotify.IN_CREATE,
		"deleted":     inotify.IN_DELETE,
		"changed":     inotify.IN_CLOSE_WRITE,
		"moved":       inotify.IN_MOVE,
		"moved_from":  inotify.IN_MOVED_FROM,
		"moved_to":    inotify.IN_MOVED_TO,
		"opened":      inotify.IN_OPEN,
		"closed":      inotify.IN_CLOSE,
		"watch_moved": inotify.IN_MOVE_SELF,
		"unmount":     inotify.IN_UNMOUNT}
	var ret uint32

	for _, flag := range conf.Options.WatchFlags {
		if flags[flag] != 0 {
			ret |= flags[flag]
		}
	}
	return ret
}
