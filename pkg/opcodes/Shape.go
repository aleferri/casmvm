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
	if s.bits >= 64 {
		return a
	}
	mask := int64(1)<<s.bits - 1
	val := a & mask
	//shapes are named i<bits> when signed and u<bits> when unsigned: a signed shape must
	//sign extend, otherwise every negative intermediate becomes a large positive and any
	//comparison against it is wrong
	if s.name[0] == 'i' && val&(int64(1)<<(s.bits-1)) != 0 {
		val |= ^mask
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
