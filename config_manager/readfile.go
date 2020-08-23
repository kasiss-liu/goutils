package configmanager

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

//ReadConfigFile is a func to load content from special file which in filepath
// to check the file type and prepare to read config file
func ReadConfigFile(filepath string, cType string) (map[string]interface{}, error) {

	_, err := os.Stat(filepath)
	if err != nil {
		return nil, errors.New("Config file is not Exist")
	}
	switch cType {
	case "ini":
		return ReadIni(filepath)
	case "json":
		return ReadJSON(filepath)
	case "yml", "yaml":
		return ReadYAML(filepath)
	default:
		return make(map[string]interface{}), errors.New("ConfigFile type is not supported")
	}

}

//ReadIni if the filetype is .ini
//to load content with this func
func ReadIni(filepath string) (map[string]interface{}, error) {

	file, error := os.Open(filepath)
	if error != nil {
		return nil, error
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var mapKey string
	deep := 0

	content := make(map[string]interface{})
	keyReg, _ := regexp.Compile("\\[(\\w+)\\]")
	varlueReg, _ := regexp.Compile("\\s*(\\w+)\\s*=\\s*(.*)")
	tmpContent := make(map[string]interface{})

	for {
		line, _, error := reader.ReadLine()
		if error != nil {
			break
		}

		lineString := string(line)
		lineString = trimCharactors(lineString)

		if lineString == "" {
			continue
		}
		commentLine, _ := regexp.MatchString("^[#|;]+.*", lineString)
		if commentLine {
			continue
		}

		matchKeys := keyReg.FindStringSubmatch(lineString)
		if len(matchKeys) > 0 {
			if deep == 0 {
				deep = 1
			} else {
				tmpContent = make(map[string]interface{})
			}
			mapKey = matchKeys[1]
			content[mapKey] = make(map[string]interface{})
			continue
		}

		matchedStrings := varlueReg.FindStringSubmatch(lineString)
		if len(matchedStrings) > 0 {
			key := matchedStrings[1]
			value := matchedStrings[2]
			if isInt, _ := regexp.MatchString("^\\d+$", value); isInt {
				intValue, _ := strconv.Atoi(value)
				if deep == 1 {
					tmpContent[key] = intValue
					content[mapKey] = tmpContent
				} else {
					content[key] = intValue
				}
			} else if isFloat, _ := regexp.MatchString("^\\d+\\.\\d+$", value); isFloat {
				flaotValue, _ := strconv.ParseFloat(value, 64)
				if deep == 1 {
					tmpContent[key] = flaotValue
					content[mapKey] = tmpContent
				} else {
					content[key] = flaotValue
				}
			} else if isArr, _ := regexp.MatchString("^\\w+,\\w+$", value); isArr {
				ss := strings.Split(value, ",")
				if len(ss) == 0 {
					continue
				}
				if deep == 1 {
					tmpContent[key] = ss
					content[mapKey] = tmpContent
				} else {
					content[key] = ss
				}

			} else {
				if deep == 1 {
					tmpContent[key] = value
					content[mapKey] = tmpContent
				} else {
					content[key] = value
				}
			}

		}

	}
	return content, nil
}

//clear unnecessary charactors
func trimCharactors(s string) string {
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r\n")
	s = strings.Trim(s, "\r")
	s = strings.Replace(s, "\"", "", -1)

	return s
}

//ReadJSON if the filetype is .json
//load content with this func
func ReadJSON(filepath string) (map[string]interface{}, error) {
	file, error := os.Open(filepath)
	if error != nil {
		return nil, error
	}
	defer file.Close()
	content, error := ioutil.ReadAll(file)
	if error != nil {
		return nil, error
	}

	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(content, &jsonMap)

	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

//ReadYAML if the filetype is .yml
//load content with this func
func ReadYAML(filepath string) (map[string]interface{}, error) {
	file, error := os.Open(filepath)
	if error != nil {
		return nil, error
	}
	defer file.Close()
	content, error := ioutil.ReadAll(file)
	if error != nil {
		return nil, error
	}

	yamlMap := make(map[string]interface{})
	err := yaml.Unmarshal(content, &yamlMap)

	if err != nil {
		return nil, err
	}
	return yamlMap, nil
}
