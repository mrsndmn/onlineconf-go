// Package onlineconf reads configuration files generated by OnlineConf.
//
// It opens indexed CDB files and maps them in the memory.
// If OnlineConf modifies them then they are automatically reopened.
package onlineconf

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/jbarham/go-cdb"
)

const configDir = "/usr/local/etc/onlineconf"

var watcher *fsnotify.Watcher

func init() {
	modules.byName = make(map[string]*Module)
	modules.byFile = make(map[string]*Module)

	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = watcher.Add(configDir)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				//log.Println("fsnotify event:", ev)

				if ev.Op&fsnotify.Create == fsnotify.Create {
					modules.Lock()
					module, ok := modules.byFile[ev.Name]
					modules.Unlock()

					if ok {
						module.reopen()
					}
				}

			case err := <-watcher.Errors:
				log.Printf("Watch %v error: %v\n", configDir, err)
			}
		}
	}()
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

type Module struct {
	m    sync.RWMutex
	name string
	file string
	cdb  *cdb.Cdb
}

func newModule(name string) *Module {
	file := fmt.Sprintf("%s/%s.cdb", configDir, name)
	cdb, err := cdb.Open(file)
	if err != nil {
		panic(err)
	}
	return &Module{name: name, file: file, cdb: cdb}
}

func (m *Module) reopen() {
	log.Printf("Reopen %s\n", m.file)
	m.m.Lock()
	defer m.m.Unlock()
	cdb, err := cdb.Open(m.file)
	if err != nil {
		log.Printf("Reopen file %v error: %v\n", m.file, err)
	} else {
		m.cdb.Close()
		m.cdb = cdb
	}
}

func (m *Module) get(path string) (byte, []byte) {
	m.m.Lock()
	defer m.m.Unlock()
	data, err := m.cdb.Data([]byte(path))
	if err != nil || len(data) == 0 {
		if err != io.EOF {
			log.Printf("Get %v:%v error: %v", m.file, path, err)
		}
		return 0, data
	}
	return data[0], data[1:]
}

// GetStringIfExists reads a string value of a named parameter from the module.
// It returns the boolean true if the parameter exists and is a string.
// In the other case it returns the boolean false and an empty string.
func (m *Module) GetStringIfExists(path string) (string, bool) {
	format, data := m.get(path)
	switch format {
	case 0:
		return "", false
	case 's':
		return string(data), true
	default:
		log.Printf("%s:%s: format is not string\n", m.name, path)
		return "", false
	}
}

// GetIntIfExists reads an integer value of a named parameter from the module.
// It returns this value and the boolean true if the parameter exists and is an integer.
// In the other case it returns the boolean false and 0.
func (m *Module) GetIntIfExists(path string) (int, bool) {
	str, ok := m.GetStringIfExists(path)
	if !ok {
		return 0, false
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("%s:%s: value is not an integer: %s\n", m.name, path, str)
		return 0, false
	}

	return i, true
}

// GetString reads a string value of a named parameter from the module.
// It returns this value if the parameter exists and is a string.
// In the other case it panics unless default value is provided in
// the second argument.
func (m *Module) GetString(path string, d ...string) string {
	if val, ok := m.GetStringIfExists(path); ok {
		return val
	} else if len(d) > 0 {
		return d[0]
	} else {
		panic(fmt.Sprintf("%s:%s key not exists and default not found", m.name, path))
	}
}

// GetInt reads an integer value of a named parameter from the module.
// It returns this value if the parameter exists and is an integer.
// In the other case it panics unless default value is provided in
// the second argument.
func (m *Module) GetInt(path string, d ...int) int {
	if val, ok := m.GetIntIfExists(path); ok {
		return val
	} else if len(d) > 0 {
		return d[0]
	} else {
		panic(fmt.Sprintf("%s:%s key not exists and default not found", m.name, path))
	}
}

var modules struct {
	sync.Mutex
	byName map[string]*Module
	byFile map[string]*Module
}

// GetModule returns a named module.
// It panics if module not exists.
func GetModule(name string) *Module {
	modules.Lock()
	defer modules.Unlock()

	if module, ok := modules.byName[name]; ok {
		return module
	}

	module := newModule(name)

	modules.byName[module.name] = module
	modules.byFile[module.file] = module

	return module
}

var tree struct {
	sync.Mutex
	module *Module
}

func getTree() *Module {
	if tree.module != nil {
		return tree.module
	}

	tree.Lock()
	defer tree.Unlock()

	if tree.module != nil {
		return tree.module
	}

	tree.module = GetModule("TREE")
	return tree.module
}

// GetStringIfExists reads a string value of a named parameter from the module "TREE".
// It returns the boolean true if the parameter exists and is a string.
// In the other case it returns the boolean false and an empty string.
func GetStringIfExists(path string) (string, bool) {
	return getTree().GetStringIfExists(path)
}

// GetIntIfExists reads an integer value of a named parameter from the module "TREE".
// It returns this value and the boolean true if the parameter exists and is an integer.
// In the other case it returns the boolean false and 0.
func GetIntIfExists(path string) (int, bool) {
	return getTree().GetIntIfExists(path)
}

// GetString reads a string value of a named parameter from the module "TREE".
// It returns this value if the parameter exists and is a string.
// In the other case it panics unless default value is provided in
// the second argument.
func GetString(path string, d ...string) string {
	return getTree().GetString(path, d...)
}

// GetInt reads an integer value of a named parameter from the module "TREE".
// It returns this value if the parameter exists and is an integer.
// In the other case it panics unless default value is provided in
// the second argument.
func GetInt(path string, d ...int) int {
	return getTree().GetInt(path, d...)
}
