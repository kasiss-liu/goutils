package Bytes2Binary

const (
	zero  = byte('0')
	one   = byte('1')
	space = byte(' ')
)

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
