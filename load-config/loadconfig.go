package loadConfig

import (
	"errors"
	"path"
	"strings"
)

// Config base info
// Config name filepath type and content
type Config struct {
	name    string
	path    string
	cType   string
	content interface{}
}

//Get key return a struct of Config
//to Get Value by a Key
func (c *Config) Get(k interface{}) *Config {
	if key, ok := k.(string); ok {
		if m, ok := (c.content).(map[string]interface{}); ok {
			if val, ok := m[key]; ok {
				if value, ok := val.(map[string]interface{}); ok {
					return &Config{c.name, c.path, c.cType, value}
				}
				return &Config{c.name, c.path, c.cType, val}
			}
		}
	}
	if m, ok := (c.content).(map[interface{}]interface{}); ok {
		if val, ok := m[k]; ok {
			if value, ok := val.(map[interface{}]interface{}); ok {
				return &Config{c.name, c.path, c.cType, value}
			}
			return &Config{c.name, c.path, c.cType, val}
		}
	}
	return &Config{}
}

//GetAll Data
func (c *Config) GetAll() interface{} {
	return c.content
}

//Int to Get an Int value
// if the type not int error will be returned
func (c *Config) Int() (int, error) {
	i, err := c.Interface()
	if err == nil {
		if f, ok := i.(int); ok {
			return f, nil
		}
	}
	if m, ok := (c.content).(int); ok {
		return m, nil
	}
	return 0, errors.New("value is not int type")
}

//String to Get a string value
// if the type not string error will be returned
func (c *Config) String() (string, error) {
	i, err := c.Interface()
	if err == nil {
		if f, ok := i.(string); ok {
			return f, nil
		}
	}
	if m, ok := (c.content).(string); ok {
		return m, nil
	}
	return "", errors.New("value is not string type")
}

//Float32 to Get a Float32 value
//if the type not float32 error will be returned
func (c *Config) Float32() (float32, error) {
	i, err := c.Interface()
	if err == nil {
		if f, ok := i.(float32); ok {
			return f, nil
		}
	}
	if m, ok := (c.content).(float32); ok {
		return m, nil
	}
	return 0, errors.New("value is not float32 type")
}

//Float64 to Get a Float64 value
//if the type not float64 error will be returned
func (c *Config) Float64() (float64, error) {
	i, err := c.Interface()
	if err == nil {
		if f, ok := i.(float64); ok {
			return f, nil
		}
	}
	if m, ok := (c.content).(float64); ok {
		return m, nil
	}
	return 0, errors.New("value is not float64 type")
}

//ArrayString to return a []string slice
// if the type wrong error will be returned
func (c *Config) ArrayString() ([]string, error) {
	//if Json or Yaml the data would be type []interface{}
	if ifv, ok := (c.content).([]interface{}); ok {
		var ss = make([]string, 0, 10)
		for _, is := range ifv {
			if s, ok := is.(string); ok {
				ss = append(ss, s)
			}
		}
		return ss, nil
	}

	if m, ok := (c.content).([]string); ok {
		return m, nil
	}
	return make([]string, 0, 0), errors.New("value is not stringArray type")
}
func (c *Config) Interface() (interface{}, error) {
	if m, ok := (c.content).(interface{}); ok {
		return m, nil
	}
	return nil, errors.New("value is not interface{} type")
}

//New to make a new Confit struct and returned
func New(name string, filepath string) *Config {
	cType := path.Ext(filepath)
	cType = strings.Trim(cType, ".")
	content, error := ReadConfigFile(filepath, cType)
	if error != nil {
		return &Config{}
	}
	return &Config{name, filepath, cType, content}
}
