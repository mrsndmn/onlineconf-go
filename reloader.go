package onlineconf

import (
	"fmt"
	"log"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
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
	module         *Mod
	modMu          *sync.RWMutex // module mutex
	ops            *ReloaderOptions
	inotifyWatcher *fsnotify.Watcher
	watherStop     chan struct{}
}

// Module returns current copy of Module by name in default onlineconf module directory.
// This module will never be updated
func Module(moduleName string) (*Mod, error) {
	mr, err := NewModuleReloader(&ReloaderOptions{Name: moduleName})
	if err != nil {
		return nil, err
	}
	return mr.Module(), err
}

// MustModule returns Module by name in default onlineconf module directory
// panics on error
func MustModule(moduleName string) *Mod {
	m, err := Module(moduleName)
	if err != nil {
		panic(err)
	}
	return m
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
		ops:        ops,
		modMu:      &sync.RWMutex{},
		watherStop: make(chan struct{}),
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

	defer func() {
		mr.inotifyWatcher = nil
	}()

	mr.watherStop <- struct{}{}

	return mr.inotifyWatcher.Close()
}

// Module returns the last successfully updated version of module
func (mr *ModuleReloader) Module() *Mod {
	mr.modMu.RLock()
	defer mr.modMu.RUnlock()
	return mr.module
}

func (mr *ModuleReloader) startWatcher() error {
	var watcher *fsnotify.Watcher

	if mr.inotifyWatcher != nil {
		return fmt.Errorf("inotify watcher is already started")
	}

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
			case <-mr.watherStop:
				log.Println("Stopping inotify watcher")
				return
			}
		}
	}()

	return nil
}

func (mr *ModuleReloader) reload() error {
	module, err := loadModuleFromFile(mr.ops.FilePath)
	if err != nil {
		log.Printf("Cant reload module %s: %#v", mr.ops.Name, err)
		return errors.Wrap(err, "cant reload module")
	}

	mr.modMu.Lock()
	mr.module = module
	mr.modMu.Unlock()
	return nil
}
