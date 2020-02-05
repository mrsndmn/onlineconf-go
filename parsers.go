package onlineconf

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

func (m *Mod) parseSimpleParams(keyStr, valStr string) {
	m.StringParams[keyStr] = valStr
	// log.Printf("str param: %s %s", keyStr, valStr)

	if intParam, err := strconv.Atoi(valStr); err == nil {
		m.IntParams[keyStr] = intParam
		// log.Printf("int param: %s %d", keyStr, intParam)
	}
	return
}

func (m *Mod) parseJSONParams(keyStr, valStr string) error {
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
