package loadConfig

import (
	"errors"
	"path"
	"strings"
)

type Config struct {
	name    string
	path    string
	cType   string
	content interface{}
}

func (c *Config) Get(key string) *Config {

	if m, ok := (c.content).(map[string]interface{}); ok {
		if val, ok := m[key]; ok {
			if value, ok := val.(map[string]interface{}); ok {
				return &Config{c.name, c.path, c.cType, value}
			}
			return &Config{c.name, c.path, c.cType, val}
		}
	}
	return &Config{c.name, c.path, c.cType, nil}
}

func (c *Config) Int() (int, error) {
	if m, ok := (c.content).(int); ok {
		return m, nil
	}
	return 0, errors.New("value is not int type !")
}

func (c *Config) String() (string, error) {
	if m, ok := (c.content).(string); ok {
		return m, nil
	}
	return "", errors.New("value is not string type !")
}

func (c *Config) Float32() (float32, error) {
	if m, ok := (c.content).(float32); ok {
		return m, nil
	}
	return 0, errors.New("value is not float32 type !")
}

func (c *Config) Float64() (float64, error) {
	if m, ok := (c.content).(float64); ok {
		return m, nil
	}
	return 0, errors.New("value is not float64 type !")
}
func (c *Config) ArrayString() ([]string, error) {
	if m, ok := (c.content).([]string); ok {
		return m, nil
	}
	return make([]string, 0, 0), errors.New("value is not stringArray type")
}

func New(name string, filepath string) *Config {
	cType := path.Ext(filepath)
	cType = strings.Trim(cType, ".")
	content, error := ReadConfigFile(filepath, cType)
	if error != nil {
		return &Config{}
	}
	return &Config{name, filepath, cType, content}
}
