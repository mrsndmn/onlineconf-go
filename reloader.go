package onlineconf

import (
	"fmt"
	"log"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// DefaultModulesDir defines default directory for modules
const DefaultModulesDir = "/usr/local/etc/onlineconf"

// ReloaderOptions specify loader options
// You can specify either FilePath or Name + Dir.
// If you sprcified only Name, DefaultModulesDir Dir will be used
type ReloaderOptions struct {
	Name     string
	Dir      string // default in `DefaultModulesDir`
	FilePath string
}

// ModuleReloader watchers for module updates and reloads it
type ModuleReloader struct {
	module         *Module
	mLock          *sync.RWMutex
	ops            *ReloaderOptions
	inotifyWatcher *fsnotify.Watcher
}

// NewModuleReloader creates new module reloader
func NewModuleReloader(ops *ReloaderOptions) (*ModuleReloader, error) {
	if ops.FilePath == "" {
		if ops.Dir == "" {
			ops.Dir = DefaultModulesDir
		}
		fileName := fmt.Sprintf("%s.cdb", ops.Name)
		filePath := path.Join(ops.Dir, fileName)
		ops.FilePath = filePath
	}

	mr := ModuleReloader{
		ops:   ops,
		mLock: &sync.RWMutex{},
	}
	err := mr.reload()
	if err != nil {
		return nil, err
	}

	err = mr.startWatcher()
	if err != nil {
		return nil, err
	}

	return &mr, nil
}

// Close closes inofitify watcher. Module will not be updated anymore.
func (mr *ModuleReloader) Close() error {
	return mr.inotifyWatcher.Close()
}

// Module returns the last successfully updated version of module
func (mr *ModuleReloader) Module() *Module {
	mr.mLock.RLock()
	defer mr.mLock.RUnlock()
	return mr.module
}

func (mr *ModuleReloader) startWatcher() error {
	var watcher *fsnotify.Watcher

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Cant init inotify watcher: %w", err)
	}

	mr.inotifyWatcher = watcher

	err = mr.inotifyWatcher.Add(mr.ops.FilePath)
	if err != nil {
		return fmt.Errorf("Cant add inotify watcher for module %s: %w", mr.ops.Name, err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Create == fsnotify.Create {
					mr.reload()
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Printf("Watch %v error: %v\n", mr.ops.Dir, err)
				}
			}
		}
	}()

	return nil
}

func (mr *ModuleReloader) reload() error {
	module, err := loadModuleFromFile(mr.ops.FilePath)
	if err != nil {
		log.Printf("Cant reload module %s: %#v", mr.ops.Name, err)
		return err
	}

	mr.mLock.Lock()
	mr.module = module
	mr.mLock.Unlock()
	return nil
}
