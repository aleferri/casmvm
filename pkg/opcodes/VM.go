package opcodes

import "github.com/aleferri/casmvm/pkg/vmio"

//VM interface on opcode execution
type VM interface {
	Frame() LocalFrame
	Goto(disp int32)
	Enter(callable int32, vals ...uint16) (LocalFrame, VMError)
	Invoke(callable int32, vals ...uint16) (LocalFrame, VMError)
	Start(callable int32, frame LocalFrame) VMError
	Leave()
	WrapError(e error) VMError
	Halt()
	Pointer() uint32
	Logger() vmio.VMLogger
	Dump(callable int32)
	DumpAll()
}
