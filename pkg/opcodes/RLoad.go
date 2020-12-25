package opcodes

import "strconv"

//RLoad load the referenced return stack offset into the evaluation stack
type RLoad struct {
	ref uint32
}

func (op *RLoad) String() string {
	return "rload " + strconv.FormatUint(uint64(op.ref), 10)
}

func (op *RLoad) Apply(vm VM) VMError {
	vm.EvalStack().Push(vm.RetStack().Load(op.ref))
	return nil
}

//MakeRLoad make an opcode of reference load
func MakeRLoad(ref uint32) Opcode {
	return &RLoad{ref}
}
