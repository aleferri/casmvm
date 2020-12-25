package opcodes

import "strconv"

type Call struct {
	offset int32
}

func (op *Call) String() string {
	return "call " + strconv.FormatInt(int64(op.offset), 10)
}

func (op *Call) Apply(vm VM) VMError {
	vm.RetStack().Push(int64(vm.Pointer()))
	vm.GotoOffset(op.offset)
	for !vm.EvalStack().Empty() {
		vm.RetStack().Push(vm.EvalStack().Pop())
	}
	return nil
}

//MakeCall instruction
func MakeCall(offset int32) Opcode {
	return &Call{offset}
}
