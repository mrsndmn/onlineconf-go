package onlineconf

// todo копировать параметры-мапки перед тем, как отдать их наружу

import "fmt"

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

// MustString returns value of a named parameter from the module.
// It panics if no such parameter or this parameter is not a string.
func (m *Module) MustString(path string) string {
	param, ok := m.StringParams[path]
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf or cant parse it %s", path))
	}
	return param
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

// MustInt returns value of a named parameter from the module.
// It panics if no such parameter or this parameter is not an int
func (m *Module) MustInt(path string) int {
	param, ok := m.IntParams[path]
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf or cant parse it %s", path))
	}
	return param
}

// MapInterfaceInterface
// Interfaces will not be copied!
func (m *Module) MapInterfaceInterface(path string) (map[interface{}]interface{}, bool) {
	param, ok := m.MapInterfaceInterfaceParams[path]
	if ok {
		clone := make(map[interface{}]interface{}, len(param))
		for k, v := range param {
			clone[k] = v
		}
		return clone, ok
	}
	return param, ok
}

// MapInterfaceInterfaceWithDef default valur will not be copied!
func (m *Module) MapInterfaceInterfaceWithDef(path string, defaultValue map[interface{}]interface{}) (map[interface{}]interface{}, bool) {
	param, ok := m.MapInterfaceInterface(path)
	if !ok {
		return defaultValue, ok
	}
	return param, ok
}

func (m *Module) MustMapInterfaceInterface(path string) map[interface{}]interface{} {
	param, ok := m.MapInterfaceInterface(path)
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf or cant parse it %s", path))
	}
	return param
}

// MapIntInt
//
//
func (m *Module) MapIntInt(path string) (map[int]int, bool) {
	param, ok := m.MapInterfaceInterface(path)
	if ok {
		clone := make(map[int]int, len(param))
		for k, v := range param {
			clone[k] = v
		}
		return clone, ok
	}
	return param, ok
}

// MapIntIntWithDef default valur will not be copied!
func (m *Module) MapIntIntWithDef(path string, defaultValue map[int]int) (map[int]int, bool) {
	param, ok := m.MapIntInt(path)
	if !ok {
		return defaultValue, ok
	}
	return param, ok
}

func (m *Module) MustMapIntInt(path string) map[int]int {
	param, ok := m.MapIntInt(path)
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf or cant parse it %s", path))
	}
	return param
}

// MapIntString
//
//
func (m *Module) MapIntString(path string) (map[int]string, bool) {
	param, ok := m.MapIntStringParams[path]
	if ok {
		clone := make(map[int]string, len(param))
		for k, v := range param {
			clone[k] = v
		}
		return clone, ok
	}
	return param, ok
}

// MapIntStringWithDef default valur will not be copied!
func (m *Module) MapIntStringWithDef(path string, defaultValue map[int]string) (map[int]string, bool) {
	param, ok := m.MapIntString(path)
	if !ok {
		return defaultValue, ok
	}
	return param, ok
}

func (m *Module) MustMapIntString(path string) map[int]string {
	param, ok := m.MapIntString(path)
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf or cant parse it %s", path))
	}
	return param
}

// MapStringInt
//
//
func (m *Module) MapStringInt(path string) (map[string]int, bool) {
	param, ok := m.MapStringIntParams[path]
	if ok {
		clone := make(map[string]int, len(param))
		for k, v := range param {
			clone[k] = v
		}
		return clone, ok
	}
	return param, ok
}

// MapStringIntWithDef default valur will not be copied!
func (m *Module) MapStringIntWithDef(path string, defaultValue map[string]int) (map[string]int, bool) {
	param, ok := m.MapStringInt(path)
	if !ok {
		return defaultValue, ok
	}
	return param, ok
}

//
func (m *Module) MustMapStringInt(path string) map[string]int {
	param, ok := m.MapStringInt(path)
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf or cant parse it %s", path))
	}
	return param
}

// MapStringString
//
//
func (m *Module) MapStringString(path string) (map[string]string, bool) {
	param, ok := m.MapStringStringParams[path]
	if ok {
		clone := make(map[string]string, len(param))
		for k, v := range param {
			clone[k] = v
		}
		return clone, ok
	}
	return param, ok
}

// MapStringStringWithDef default valur will not be copied!
func (m *Module) MapStringStringWithDef(path string, defaultValue map[string]string) (map[string]string, bool) {
	param, ok := m.MapStringString(path)
	if !ok {
		return defaultValue, ok
	}
	return param, ok
}

func (m *Module) MustMapStringString(path string) map[string]string {
	param, ok := m.MapStringString(path)
	if !ok {
		panic(fmt.Errorf("Missing required parameter in onlineconf or cant parse it %s", path))
	}
	return param
}
