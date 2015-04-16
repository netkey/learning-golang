package dir

import (
	"golang.org/x/exp/inotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Dir struct {
	path        string
	name        string
	flags       uint32
	files       map[string]*os.FileInfo
	directories map[string]*Dir
	recurse     bool
	Watcher     *inotify.Watcher
}

func NewDir(path string, name string, flags uint32, recurse bool) *Dir {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return nil
	}
	return &Dir{
		path:        filepath.Clean(path),
		name:        name,
		flags:       flags,
		files:       nil,
		directories: nil,
		recurse:     recurse,
		Watcher:     watcher,
	}
}

func (dir *Dir) GetPath() string {
	return dir.path
}

func (dir *Dir) GetName() string {
	return dir.name
}

func (dir *Dir) GetFlags() uint32 {
	return dir.flags
}

func (dir *Dir) GetFiles() map[string]*os.FileInfo {
	return dir.files
}

func (dir *Dir) GetChildren() (map[string]*os.FileInfo, map[string]*Dir) {
	return dir.files, dir.directories
}

func (dir *Dir) GetDirectories() map[string]*Dir {
	return dir.directories
}

func (dir *Dir) AnyChildDirs() int {
	return len(dir.directories)
}

func (dir *Dir) GetChildDir(name string) *Dir {
	if name == dir.GetPath() {
		return dir
	}

	children := dir.GetDirectories()
	for _, child := range children {
		if child.GetPath() == name {
			return child
		} else if ret := child.GetChildDir(name); ret != nil {
			return ret
		}
	}
	return nil
}

func (dir *Dir) GetChildFile(name string) *os.FileInfo {
	return dir.files[name]
}

func (dir *Dir) enumerateChildren() error {
	log.Println("Enumerating children of watch directory:", dir.name)

	list, err := ioutil.ReadDir(dir.path)
	if err != nil {
		return err
	}

	dir.files = make(map[string]*os.FileInfo)
	dir.directories = make(map[string]*Dir)
	for _, file := range list {
		//index by path instead of name
		path := dir.path + "/" + file.Name()

		//If it's a directory we want to create a whole new Dir struct
		//and index it in the parent directory's list of subdirectories
		if file.IsDir() == true {
			log.Println(dir.name, "-- DIR:", path)

			//children inherit the recursive flag and watch flags from their parent
			dir.directories[path] = NewDir(path, file.Name(), dir.flags, dir.recurse)
			if dir.recurse == true {
				child := dir.directories[path]
				if err := child.enumerateChildren(); err != nil {
					log.Println(err)
				}
			}
		} else {
			//If it's a file we just get a pointer to
			//the fileinfo struct, again indexing by path
			log.Println(dir.name, "-- FILE:", file.Name())
			dir.files[path] = &file
		}
	}
	return nil
}

//Meant to be used when starting a watch
func (dir *Dir) watchChildren() error {
	for _, child := range dir.GetDirectories() {
		if child.recurse == true {
			log.Println("Adding", child.path, "to watch")

			if err := dir.Watcher.AddWatch(child.path, child.flags); err != nil {
				return err
			}

			if err := child.watchChildren(); err != nil {
				return err
			}

		} else {
			log.Println("Child", child.path, "not recursive")
		}
	}
	return nil
}

func (dir *Dir) StartWatch() error {
	if dir.files == nil && dir.directories == nil {
		if err := dir.enumerateChildren(); err != nil {
			return err
		}
	}
	if err := dir.Watcher.AddWatch(dir.path, dir.flags); err != nil {
		return err
	}
	return dir.watchChildren()
}
