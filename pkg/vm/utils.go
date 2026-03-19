package vm

func convToUint16(vals [2]uint8) uint16 {
	return uint16(vals[0]) | uint16(vals[1])<<8
}

func convToUint8(val uint16) []uint8 {
	return []uint8{uint8(val), uint8(val >> 8)}
}
