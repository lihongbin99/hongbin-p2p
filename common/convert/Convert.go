package convert

const (
	length = 4
)

func I2b(num int32) []byte {
	return []byte{uint8(num >> 24), uint8(num >> 16), uint8(num >> 8), uint8(num >> 0)}
}

func B2i(b []byte) int32 {
	return int32(b[3]) + int32(b[2])<<8 + int32(b[1])<<16 + int32(b[0])<<24
}
