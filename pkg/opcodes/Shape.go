package opcodes

//Shape of a value
type Shape struct {
	name string
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

var ShortShape = Shape{"i16", 0}
var UShortShape = Shape{"u16", 0}
var IntShape = Shape{"i32", 0}
var UIntShape = Shape{"u32", 0}
