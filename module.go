package onlineconf

import (
	"context"
	"fmt"
	"io"

	"github.com/alldroll/cdb"
	"github.com/pkg/errors"
)

// ErrInvalidCDB means that cdb is invalid
var ErrInvalidCDB = errors.New("cdb is inconsistent")

// Mod is a structure that associated with onlineconf module file.
type Mod struct {
	StringParams map[string]string
	IntParams    map[string]int

	RawJSONParams         map[string]string // Here will be all JSON params (not parsed)
	JSONObjectParams      map[string]map[string]interface{}
	MapIntIntParams       map[string]map[int]int
	MapIntStringParams    map[string]map[int]string
	MapStringIntParams    map[string]map[string]int
	MapStringStringParams map[string]map[string]string
	JSONArrayParams       map[string][]interface{}
	SliceStringParams     map[string][]string
	SliceIntParams        map[string][]int
}

func newEmptyModule() *Mod {
	return &Mod{
		StringParams: map[string]string{},
		IntParams:    map[string]int{},

		RawJSONParams:         map[string]string{},
		JSONObjectParams:      map[string]map[string]interface{}{},
		MapIntIntParams:       map[string]map[int]int{},
		MapIntStringParams:    map[string]map[int]string{},
		MapStringIntParams:    map[string]map[string]int{},
		MapStringStringParams: map[string]map[string]string{},
		JSONArrayParams:       map[string][]interface{}{},
		SliceStringParams:     map[string][]string{},
		SliceIntParams:        map[string][]int{},
	}
}

// NewModule parses cdb file and copies all content to local maps.
// Module returned by this method will never be updated
func NewModule(reader io.ReaderAt) (*Mod, error) {

	cdbReader, err := cdb.New().GetReader(reader)
	if err != nil {
		return nil, fmt.Errorf("Cant cant cdb reader for module: %w", err)
	}

	module := newEmptyModule()

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

func (m *Mod) fillParams(cdb cdb.Reader) error {
	cdbIter, err := cdb.Iterator()
	if err != nil {
		return errors.Wrap(err, "cant get cdb iterator")
	}

	for {
		record := cdbIter.Record()
		if record == nil {
			break
		}

		key, err := record.KeyBytes()
		if err != nil {
			return errors.Wrap(err, "cant read cdb key")
		}

		val, err := record.ValueBytes()
		if err != nil {
			return errors.Wrap(err, "cant read cdb value")
		}

		if len(val) < 1 {
			return fmt.Errorf("Onlineconf value must contain at least 1 byte: `typeByte|ParamData`")
		}

		// log.Printf("oc parsing: %s %s", string(key), string(val))

		// val's first byte defines datatype of config value
		// onlineconf currently knows 's' and 'j' data types
		paramTypeByte := val[0]
		keyStr := string(key)
		valStr := string(val[1:])
		if paramTypeByte == 's' { // params type string
			m.parseSimpleParams(keyStr, valStr)
		} else if paramTypeByte == 'j' { // params type JSON

			// todo сделать парсинг json параметров опциональным
			err := m.parseJSONObjectParams(keyStr, valStr)
			if err != nil {
				arrayerr := m.parseJSONArrayParams(keyStr, valStr)
				if arrayerr != nil {
					return fmt.Errorf("Cant parse json. It is neither array nor json object: %#v %#v", err, arrayerr)
				}

			}
		} else {
			return fmt.Errorf("Unknown paramTypeByte: %#v for key %s", paramTypeByte, keyStr)
		}

		if !cdbIter.HasNext() {
			break
		}

		_, err = cdbIter.Next()
		if err != nil {
			return errors.Wrap(err, "cant get next cdb record")
		}
	}

	return nil
}

type ctxConfigModuleKey struct{}

// WithContext returns a new Context that carries value module
func (m *Mod) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxConfigModuleKey{}, m)
}

// ModuleFromContext retrieves a config module from context.
// Returns empty module if context has no ctxConfigModuleKey
func ModuleFromContext(ctx context.Context) *Mod {
	m, ok := ctx.Value(ctxConfigModuleKey{}).(*Mod)
	if !ok {
		return newEmptyModule()
	}
	return m
}
