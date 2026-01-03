package opcodes

import "github.com/aleferri/casmvm/pkg/vmio"

//VM interface on opcode execution
type VM interface {
	Frame() LocalFrame
	Pointer() uint32
	Goto(disp int32)
	Enter(callable int32, vals ...uint16) (LocalFrame, VMError)
	Invoke(callable int32, vals ...uint16) (LocalFrame, VMError)
	Start(callable int32, frame LocalFrame) VMError
	Leave()
	Halt()
	WordSize() int
	WrapError(e error) VMError
	Logger() vmio.VMLogger
	Dump(callable int32)
	DumpAll()
}
