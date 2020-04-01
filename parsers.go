package onlineconf

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func (m *Mod) parseJSONObjectParams(keyStr, valStr string) error {
	m.RawJSONParams[keyStr] = valStr

	byteVal := []byte(valStr)

	mapStringInterface := make(map[string]interface{})
	err := json.Unmarshal(byteVal, &mapStringInterface)
	if err != nil {
		return nil
	}
	m.JSONObjectParams[keyStr] = mapStringInterface

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

func (m *Mod) parseJSONArrayParams(keyStr, valStr string) error {
	m.RawJSONParams[keyStr] = valStr

	fmt.Printf("valstr %s", valStr)
	byteVal := []byte(valStr)

	mapInterfaceSlice := make([]interface{}, 0)
	err := json.Unmarshal(byteVal, &mapInterfaceSlice)
	if err != nil {
		return nil
	}
	m.JSONArrayParams[keyStr] = mapInterfaceSlice

	strSlice := make([]string, 0)
	err = json.Unmarshal(byteVal, &strSlice)
	if err != nil {
		return nil
	}
	m.SliceStringParams[keyStr] = strSlice

	intSlice := make([]int, 0)
	err = json.Unmarshal(byteVal, &intSlice)
	if err != nil {
		return nil
	}
	m.SliceIntParams[keyStr] = intSlice

	return nil
}
