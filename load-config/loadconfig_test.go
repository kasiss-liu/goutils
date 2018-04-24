package loadConfig

import (
	"fmt"
	"testing"
)

func TestLoadConfigIni(t *testing.T) {

	config := New("test", "./test/conf.ini")

	homePage, err := config.Get("homePage").String()
	if err == nil {
		fmt.Println(homePage)
	} else {
		fmt.Println(err.Error())
	}

	name, err := config.Get("author").Get("name").String()
	if err == nil {
		fmt.Println(name)
	} else {
		fmt.Println(err.Error())
	}

	age, err := config.Get("author").Get("age").Int()
	if err == nil {
		fmt.Println(age)
	} else {
		fmt.Println(err.Error())
	}

	// Output:
	//	https://github.com/kasiss-liu
	//	kasiss
	//	28
}

func LoadConfigJson(t *testing.T) {

	config := New("test", "./test/conf.json")
	var err error
	homePage, err := config.Get("homePage").String()
	if err == nil {
		fmt.Println(homePage)
	} else {
		fmt.Println(err.Error())
	}

	name, err := config.Get("author").Get("name").String()
	if err == nil {
		fmt.Println(name)
	} else {
		fmt.Println(err.Error())
	}

	age, err := config.Get("author").Get("age").Float64()
	if err == nil {
		fmt.Println(age)
	} else {
		fmt.Println(err.Error())
	}

	// Output:
	//	https://github.com/kasiss-liu
	//	kasiss
	//	28
}

func TestLoadConfigYml(t *testing.T) {

	config := New("test", "./test/conf.yml")

	homePage, err := config.Get("homePage").String()

	if err == nil {
		fmt.Println(homePage)
	} else {
		fmt.Println(err.Error())
	}

	name, err := config.Get("author").Get("name").String()
	if err == nil {
		fmt.Println(name)
	} else {
		fmt.Println(err.Error())
	}

	age, err := config.Get("author").Get("age").Int()
	if err == nil {
		fmt.Println(age)
	} else {
		fmt.Println(err.Error())
	}

	// Output:
	//	https://github.com/kasiss-liu
	//	kasiss
	//	28
}
