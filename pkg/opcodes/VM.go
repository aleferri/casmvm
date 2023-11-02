package opcodes

import "github.com/aleferri/casmvm/pkg/vmio"

//VM interface on opcode execution
type VM interface {
	Frame() LocalFrame
	Goto(disp int32)
	Enter(frame int32, vals ...uint16) (LocalFrame, VMError)
	Leave()
	Invoke(fIndex int32, frame LocalFrame) VMError
	WrapError(e error) VMError
	Halt()
	Pointer() uint32
	Logger() vmio.VMLogger
	Dump(frame int32)
	DumpAll()
}
