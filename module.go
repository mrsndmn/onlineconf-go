package onlineconf

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/alldroll/cdb"
	"github.com/pkg/errors"
)

// ErrInvalidCDB means that cdb is invalid
var ErrInvalidCDB = errors.New("cdb is inconsistent")

// Module is a structure that associated with onlineconf module file.
type Module struct {
	StringParams map[string]string
	IntParams    map[string]int

	RawJSONParams               map[string]string // Here will be all JSON params (not parsed)
	MapInterfaceInterfaceParams map[string]map[interface{}]interface{}
	MapIntIntParams             map[string]map[int]int
	MapIntStringParams          map[string]map[int]string
	MapStringIntParams          map[string]map[string]int
	MapStringStringParams       map[string]map[string]string
}

// NewModule parses cdb file and copies all content to local maps
func NewModule(reader io.ReaderAt) (*Module, error) {

	cdbReader, err := cdb.New().GetReader(reader)
	if err != nil {
		return nil, fmt.Errorf("Cant cant cdb reader for module: %w", err)
	}

	module := &Module{
		StringParams: map[string]string{},
		IntParams:    map[string]int{},

		RawJSONParams:               map[string]string{},
		MapInterfaceInterfaceParams: map[string]map[interface{}]interface{}{},
		MapIntIntParams:             map[string]map[int]int{},
		MapIntStringParams:          map[string]map[int]string{},
		MapStringIntParams:          map[string]map[string]int{},
		MapStringStringParams:       map[string]map[string]string{},
	}

	// todo подумать, как будут обновляться модули
	// кажется, что горутинка при обновлении файлика должна
	// генерить новый модуль и отдавать ссылку нна него по запросу
	// пока файлик не обновится еще раз
	err = module.fillParams(cdbReader)
	if err != nil {
		return nil, err
	}
	return module, nil
}

func (m *Module) fillParams(cdb cdb.Reader) error {
	cdbIter, err := cdb.Iterator()
	if err != nil {
		return errors.Wrap(err, "cant get cdb iterator")
	}

	record := cdbIter.Record()
	_, ks := record.Key()
	log.Printf("1 rec: %d", ks)

	for {
		record := cdbIter.Record()
		if record == nil {
			break
		}

		keyReader, keySize := record.Key()
		key := make([]byte, int(keySize))
		if _, err = keyReader.Read(key); err != nil {
			return errors.Wrap(err, "cant read cdb key")
		}

		valReader, valSize := record.Value()
		val := make([]byte, int(valSize))
		if _, err = valReader.Read(val); err != nil {
			return errors.Wrap(err, "cant read cdb value")
		}

		if len(val) < 1 {
			return fmt.Errorf("Onlineconf value must contain at least 1 byte: `typeByte|ParamData`")
		}

		log.Printf("oc parsing: %s %s", string(key), string(val))

		// val's first byte defines datatype of config value
		// onlineconf currently knows 's' and 'j' data types
		paramTypeByte := val[0]
		keyStr := string(key)
		valStr := string(val[1:])
		if paramTypeByte == 's' { // params type string
			m.parseSimpleParams(keyStr, valStr)
		} else if paramTypeByte == 'j' { // params type JSON
			err := m.parseJSONParams(keyStr, valStr)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unknown paramTypeByte: %#v for key %s", paramTypeByte, keyStr)
		}

		if !cdbIter.HasNext() {
			break
		}

		_, err := cdbIter.Next()
		if err != nil {
			return errors.Wrap(err, "cant get next cdb record")
		}
	}

	return nil
}

func (m *Module) parseSimpleParams(keyStr, valStr string) {
	m.StringParams[keyStr] = valStr
	log.Printf("str param: %s %s", keyStr, valStr)

	if intParam, err := strconv.Atoi(valStr); err == nil {
		m.IntParams[keyStr] = intParam
		log.Printf("int param: %s %d", keyStr, intParam)
	}
	return
}

func (m *Module) parseJSONParams(keyStr, valStr string) error {
	m.RawJSONParams[keyStr] = valStr

	byteVal := []byte(valStr)

	mapInterfaceInterface := make(map[interface{}]interface{})
	err := json.Unmarshal(byteVal, &mapInterfaceInterface)
	if err != nil {
		return errors.Wrapf(err, "invalid json in parameter %s", keyStr)
	}
	m.MapInterfaceInterfaceParams[keyStr] = mapInterfaceInterface

	mapStrStr := make(map[string]string)
	err = json.Unmarshal(byteVal, &mapStrStr)
	if err != nil {
		return nil
	}

	mapStrInt := make(map[string]int)
	mapIntStr := make(map[int]string)
	mapIntInt := make(map[int]int)

	for k, v := range mapStrStr {
		var intK, intV int
		intK, keyErr := strconv.Atoi(k)
		intV, valErr := strconv.Atoi(v)

		if keyErr == nil {
			mapIntStr[intK] = v
		}
		if valErr == nil {
			mapStrInt[k] = intV
		}
		if valErr == nil && keyErr == nil {
			mapIntInt[intK] = intV
		}
	}

	m.MapIntIntParams[keyStr] = mapIntInt
	m.MapIntStringParams[keyStr] = mapIntStr
	m.MapStringIntParams[keyStr] = mapStrInt
	m.MapStringStringParams[keyStr] = mapStrStr
	return nil
}

// String returns value of a named parameter from the module.
// It returns the boolean true if the parameter exists and is a string.
// In the other case it returns the boolean false and an empty string.
func (m *Module) String(path string) (string, bool) {
	param, ok := m.StringParams[path]
	return param, ok
}

// StringWithDef returns value of a named parameter from the module.
// It returns the boolean true if the parameter exists and is a string.
// In the other case it returns the boolean false and an empty string.
func (m *Module) StringWithDef(path string, defaultValue string) (string, bool) {
	param, ok := m.StringParams[path]
	if !ok {
		return defaultValue, ok
	}
	return param, ok
}

// Int returns value of a named parameter from the module.
// It returns the boolean true if the parameter exists and is an int.
// In the other case it returns the boolean false and zero.
func (m *Module) Int(path string) (int, bool) {
	param, ok := m.IntParams[path]
	return param, ok
}

// IntWithDef returns value of a named parameter from the module.
// It returns the boolean true if the parameter exists and is an int.
// In the other case it returns the boolean false and zero.
func (m *Module) IntWithDef(path string, defaultValue int) (int, bool) {
	param, ok := m.IntParams[path]
	if !ok {
		return defaultValue, ok
	}
	return param, ok
}

// MustString returns value of a named parameter from the module.
// It panics if no such parameter or this parameter is not a string.
func (m *Module) MustString(path string) string {
	param, ok := m.StringParams[path]
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf %s", path))
	}
	return param
}

// MustInt returns value of a named parameter from the module.
// It panics if no such parameter or this parameter is not an int
func (m *Module) MustInt(path string) int {
	param, ok := m.IntParams[path]
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf %s", path))
	}
	return param
}
