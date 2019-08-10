package beacon

// source https://github.com/suapapa/go_eddystone/blob/master/cast.go
func fixTofloat32(a uint16) float32 {
	if a&0x8000 == 0 {
		return float32(a) / 256.0
	}
	return -(float32(^a) + 1) / 256.0
}

func bytesToUint16(a []byte) (v uint16) {
	_ = a[1]
	v = uint16(a[0])<<8 | uint16(a[1])
	return
}

func byteToInt(a byte) (v int) {
	v = int(a)
	if v&0x80 != 0 {
		v = -((^v + 1) & 0xff)
	}
	return
}
