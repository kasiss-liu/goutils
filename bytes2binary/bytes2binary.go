//本类包功能为将字节与二进制字符串互转
//参考借鉴了互联网内容
//https://blog.csdn.net/qq245671051/article/details/52693359

package Bytes2Binary

const (
	zero  = byte('0')
	one   = byte('1')
	space = byte(' ')
)

var uint8arr [8]uint8

func init() {
	uint8arr[0] = 128
	uint8arr[1] = 64
	uint8arr[2] = 32
	uint8arr[3] = 16
	uint8arr[4] = 8
	uint8arr[5] = 4
	uint8arr[6] = 2
	uint8arr[7] = 1
}

//将一个byte存入字节数组
func appendByte2BinaryArr(bs []byte, b byte) []byte {
	var a byte
	for i := 0; i < 8; i++ {
		a = b
		b <<= 1
		b >>= 1
		switch a {
		case b:
			bs = append(bs, zero)
		default:
			bs = append(bs, one)
		}
		b <<= 1
	}
	return bs
}

//将单个字节转化为二进制字符串
func Byte2Binary(b byte) string {
	return string(appendByte2BinaryArr([]byte{}, b))
}

//将字节数组转化为二进制字符串
func Bytes2Binary(bs []byte) string {
	var rs []byte
	for _, v := range bs {
		rs = appendByte2BinaryArr(rs, v)
		rs = append(rs, space)
	}
	return string(rs)
}

//将二进制字符转化为 []byte
func Binary2Bytes(binary string) []byte {
	bs := make([]byte, 0, 10)
	l := len(binary)
	if l == 0 {
		return bs
	}
	mo := l % 8
	l /= 8
	if mo != 0 {
		l++
	}
	bs = make([]byte, 0, l)
	mo = 8 - mo
	var n uint8
	for i, b := range []byte(binary) {
		m := (i + mo) % 8
		switch b {
		case one:
			n += uint8arr[m]
		}
		if m == 7 {
			bs = append(bs, n)
			n = 0
		}
	}
	return bs
}
