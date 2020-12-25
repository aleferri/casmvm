package opcodes

//Stack is the standard Stack interface
type Stack interface {
	Push(value int64)
	Pop() int64
	Load(offset uint32) int64
	Store(offset uint32, value int64)
	Empty() bool
}
