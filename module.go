package onlineconf

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/alldroll/cdb"
	"github.com/pkg/errors"
)

// Module is a structure that associated with onlineconf module file.
type Module struct {
	StringParams map[string]string
	IntParams    map[string]int
}

// NewModule parses cdb file and copies all content to local maps
func NewModule(reader io.ReaderAt) (*Module, error) {

	cdbReader, err := cdb.New().GetReader(reader)
	if err != nil {
		return nil, fmt.Errorf("Cant cant cdb reader for module: %w", err)
	}

	module := &Module{}

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

// для того, чтобы не приходилось каждый раз парсить
// содержимое конфига, это можно сделать один раз.
// для этого надо знать, от онлайнконфа, какого типа даннный парамерт
// так же, как это сделано сейас, например, для JSON можно писать подсказки
// в cdb файле и парсить или не парсить данное число.
// Как вариант, можно все параметры, для которых не указан тип,
// пытаться распарсить и как число, и как строку, и как число разных типов uint64 и т д
// то, что получилось, складывать в отдельные мапки и при обращении вообще не ходить в cdb файл
// до этого, наверное, интересно побенчить, на сколько мапка будет быстрее cdb
func (m *Module) fillParams(cdb cdb.Reader) error {
	stringParams := map[string]string{}
	intParams := map[string]int{}

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
		if paramTypeByte == 's' { // params type string
			keyStr := string(key)
			valStr := string(val[1:])

			// todo? check val first byte.
			// if its 0
			stringParams[keyStr] = valStr
			log.Printf("str param: %s %s", keyStr, valStr)

			if intParam, err := strconv.Atoi(valStr); err == nil {
				intParams[keyStr] = intParam
				log.Printf("int param: %s %d", keyStr, intParam)
			}
		} else if paramTypeByte == 'j' {
			// not supported yet
			// todo support json params
			panic("not supported record type")
		} else {
			return fmt.Errorf("Unknown paramTypeByte: %#v for key %s", paramTypeByte, string(key))
		}

		if !cdbIter.HasNext() {
			break
		}

		_, err := cdbIter.Next()
		if err != nil {
			return errors.Wrap(err, "cant get next cdb record")
		}
	}

	m.IntParams = intParams
	m.StringParams = stringParams
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
