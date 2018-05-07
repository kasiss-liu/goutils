package fluentLogger

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {

	RegisterFileLogger("test", "./test", "testlogger", PARTITION_NONE)
	err := StartLogger()
	if err != nil {
		t.Log(err.Error())
	}
	for i := 0; i < 10; i++ {
		tm := time.Now()
		d := i % 4
		switch d {
		case 0:
			Log("test", fmt.Sprintf("%s%d", "testlog_", i), tm)
		case 1:
			Notice("test", fmt.Sprintf("%s%d", "testlog_", i), tm)
		case 2:
			Warning("test", fmt.Sprintf("%s%d", "testlog_", i), tm)
		case 3:
			Error("test", fmt.Sprintf("%s%d", "testlog_", i), tm)
		}

	}
	logLock.Wait()
}

func BenchmarkLogger(b *testing.B) {
	RegisterFileLogger("app", "./test", "benchmark_new", PARTITION_NONE)
	err := StartLogger()
	if err != nil {
		b.Log(err.Error())
	}

	for i := 0; i < b.N; i++ {
		Log("app", fmt.Sprintf("%s%d", "Cras eu dolor lorem. Cras justo mauris, rhoncus in mauris ac, pellentesque pulvinar metus. Suspendisse consectetur consequat diam, ac dignissim mauris gravida vitae. Vestibulum blandit vestibulum mi a viverra.", i), time.Now())
	}
	logLock.Wait()
}
