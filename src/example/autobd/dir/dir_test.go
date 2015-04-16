package dir

import (
	"github.com/SaviorPhoenix/autobd/config"
	"golang.org/x/exp/inotify"
	"testing"
)

func Test_StartWatch(t *testing.T) {
	conf := config.ParseConfig("../config.toml")
	watchFlags := conf.GetFlags()
	dir := NewDir("../test", "../test", watchFlags, conf.Options.Recursive)

	watcher, err := inotify.NewWatcher()
	if err != nil {
		t.Error("Failed to allocate inotify watcher")
	}

	if dir == nil {
		t.Error("dir is nil")
	}

	if err := dir.StartWatch(watcher); err != nil {
		t.Error(err)
	}

	if dir.GetDirectories() == nil {
		t.Error("dir.GetDirectories() is nil")
	}
	if dir.GetFiles() == nil {
		t.Error("dir.GetFiles() is nil")
	}
	directories, files := dir.GetChildren()
	if directories == nil || files == nil {
		t.Error("dir.GetChildren() is nil")
	}

	child := dir.GetChildDir("../test/two")
	if child == nil {
		t.Error("Could not get child directory from dir.GetChildDir()")
	}
}
