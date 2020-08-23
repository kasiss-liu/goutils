package configmanager

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
		t.Error(err.Error())
	}

	name, err := config.Get("author").Get("name").String()
	if err == nil {
		fmt.Println(name)
	} else {
		t.Error(err.Error())
	}

	age, err := config.Get("author").Get("age").Int()
	if err == nil {
		fmt.Println(age)
	} else {
		t.Error(err.Error())
	}

	per, err := config.Get("author").Get("percentage").Float64()
	if err == nil {
		fmt.Println(per)
	} else {
		t.Error(err.Error())
	}

	arr, err := config.Get("author").Get("keywords").ArrayString()
	if err == nil {
		fmt.Println(arr)
	} else {
		t.Error(err.Error())
	}

	// Output:
	//	https://github.com/kasiss-liu
	//	kasiss
	//	28
	//  0.1
	//	[kasiss liu]
}

func TestLoadConfigJson(t *testing.T) {

	config := New("test", "./test/conf.json")
	var err error
	homePage, err := config.Get("homePage").String()
	if err == nil {
		fmt.Println(homePage)
	} else {
		t.Error(err.Error())
	}

	name, err := config.Get("author").Get("name").String()
	if err == nil {
		fmt.Println(name)
	} else {
		t.Error(err.Error())
	}

	age, err := config.Get("author").Get("age").Float64()
	if err == nil {
		fmt.Println(age)
	} else {
		t.Error(err.Error())
	}
	per, err := config.Get("author").Get("percentage").Float64()
	if err == nil {
		fmt.Println(per)
	} else {
		t.Error(err.Error())
	}
	arr, err := config.Get("author").Get("keywords").ArrayString()
	if err == nil {
		fmt.Println(arr)
	} else {
		t.Error(err.Error())
	}

	// Output:
	//	https://github.com/kasiss-liu
	//	kasiss
	//	28
	//	0.1
	//	[kasiss liu]
}

func TestLoadConfigYml(t *testing.T) {

	config := New("test", "./test/conf.yml")

	homePage, err := config.Get("homePage").String()

	if err == nil {
		fmt.Println(homePage)
	} else {
		t.Error(err.Error())
	}

	name, err := config.Get("author").Get("name").String()
	if err == nil {
		fmt.Println(name)
	} else {
		t.Error(err.Error())
	}

	age, err := config.Get("author").Get("age").Int()
	if err == nil {
		fmt.Println(age)
	} else {
		t.Error(err.Error())
	}
	per, err := config.Get("author").Get("percentage").Float64()
	if err == nil {
		fmt.Println(per)
	} else {
		t.Error(err.Error())
	}

	arr, err := config.Get("author").Get("keywords").ArrayString()
	if err == nil {
		fmt.Println(arr)
	} else {
		t.Error(err.Error())
	}

	// Output:
	//	https://github.com/kasiss-liu
	//	kasiss
	//	28
	//	0.1
	//	[kasiss liu]
}
