package Bytes2Binary

import (
	"testing"
)

func TestByte2Binary(t *testing.T) {
	a := Byte2Binary('0')
	println(a)
	// Output: 00110000
}

func TestBytes2Binary(t *testing.T) {
	a := Bytes2Binary([]byte{'0', '1'})
	println(a)
	// Output: 00110000 00110001
}
