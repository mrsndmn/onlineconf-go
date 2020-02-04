package onlineconf

import (
	"fmt"
	"io"
	"log"

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
