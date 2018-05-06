package fluentLogger

import (
	"testing"
)

var testDir = "./test/"

func TestPartitionNoneWriter(t *testing.T) {
	fwriter := NewFileWriter("./test", "test", PARTITION_NONE)
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
	fwriter := NewFileWriter("./test", "test", PARTITION_DAY)
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
	fwriter := NewFileWriter("./test", "test", PARTITION_HOUR)
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
