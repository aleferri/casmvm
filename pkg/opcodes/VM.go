package opcodes

//VM interface on opcode execution
type VM interface {
	EvalStack() Stack
	RetStack() Stack
	Goto(ptr uint32)
	GotoOffset(disp int32)
	WrapError(e error) VMError
	Halt()
	Pointer() uint32
}
