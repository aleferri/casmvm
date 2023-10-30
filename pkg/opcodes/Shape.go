package opcodes

//Shape of a value
type Shape struct {
	name string
	bits uint16
	indx uint16
}

//Index of a shape
func (s Shape) Index() uint16 {
	return s.indx
}

//Name of a shape
func (s Shape) Name() string {
	return s.name
}

func (s Shape) Reshape(a int64) int64 {
	bitmask := 0xF
	signmask := 0x7
	for i := 4; i < int(s.bits); i += 4 {
		bitmask = bitmask<<4 + 0xF
		signmask = signmask << 4
	}

	val := a & int64(bitmask)
	sign := a & int64(signmask)
	if len(s.name) == 3 && s.name[0] == 's' {
		sext := sign
		for t := 1; t < 64-int(s.bits); t++ {
			sext = sext<<1 + sign
		}
		val = val | sext
	}
	return val
}

var ByteShape = Shape{"i8", 8, 0}
var UByteShape = Shape{"u8", 8, 0}
var ShortShape = Shape{"i16", 16, 0}
var UShortShape = Shape{"u16", 16, 0}
var IntShape = Shape{"i32", 32, 0}
var UIntShape = Shape{"u32", 32, 0}
var LongShape = Shape{"i64", 64, 0}
var ULongShape = Shape{"u64", 64, 0}
