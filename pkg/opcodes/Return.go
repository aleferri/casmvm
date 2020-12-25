package opcodes

import "strconv"

type Return struct {
	depth uint32
}

func (op *Return) String() string {
	return "ret " + strconv.FormatUint(uint64(op.depth), 10)
}

func (op *Return) Apply(vm VM) VMError {
	unwind := op.depth
	for unwind > 0 {
		unwind--
		vm.RetStack().Pop()
	}
	if vm.RetStack().Empty() {
		vm.Halt()
	} else {
		dest := vm.RetStack().Pop()
		vm.Goto(uint32(dest & 0xFFFFFFFF))
	}
	return nil
}

//MakeReturn make an opcode of reference load
func MakeReturn(depth uint32) Opcode {
	return &Return{depth}
}
