package loadConfig

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ReadConfigFile(filepath string, cType string) (map[string]interface{}, error) {

	_, error := os.Stat(filepath)
	if error != nil {
		return nil, errors.New("Config file is not Exist !")
	}

	switch cType {
	case "ini":
		return ReadIni(filepath)
	case "json":
		return ReadJson(filepath)
	default:
		return make(map[string]interface{}), errors.New("ConfigFile type is not supported !")
	}

}

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
	varlueReg, _ := regexp.Compile("(\\w+)=(.*)")
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
			if isInt, _ := regexp.MatchString("\\d+", value); isInt {
				intValue, _ := strconv.Atoi(value)
				if deep == 1 {
					tmpContent[key] = intValue
					content[mapKey] = tmpContent
				} else {
					content[key] = intValue
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

func trimCharactors(s string) string {
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r\n")
	s = strings.Trim(s, "\r")
	s = strings.Replace(s, "\"", "", -1)
	s = strings.Replace(s, " ", "", -1)

	return s
}

func ReadJson(filepath string) (map[string]interface{}, error) {
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
