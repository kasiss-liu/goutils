package fluentLogger

import (
	"testing"
)

var testDir = "./test/"

func TestPartitionNoneWriter(t *testing.T) {
	fwriter := NewFileWriter("./test", "test", PartitionNone)
	err := fwriter.Prepare()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = fwriter.Prepare()
	if err != nil {
		t.Fatal(err.Error())
	}

	n, err := fwriter.Write([]byte("hello"))
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("write bytes:", n)

	n, err = fwriter.WriteString("world")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("writestring bytes:", n)
	fwriter.Close()
}

func TestPartitionDayWriter(t *testing.T) {
	fwriter := NewFileWriter("./test", "test", PartitionDay)
	err := fwriter.Prepare()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = fwriter.Prepare()
	if err != nil {
		t.Fatal(err.Error())
	}

	n, err := fwriter.Write([]byte("hello"))
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("write bytes:", n)

	n, err = fwriter.WriteString("world")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("writestring bytes:", n)
	fwriter.Close()
}

func TestPartitionHourWriter(t *testing.T) {
	fwriter := NewFileWriter("./test", "test", PartitionHour)
	err := fwriter.Prepare()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = fwriter.Prepare()
	if err != nil {
		t.Fatal(err.Error())
	}

	n, err := fwriter.Write([]byte("hello"))
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("write bytes:", n)

	n, err = fwriter.WriteString("world")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("writestring bytes:", n)
	fwriter.Close()
}
