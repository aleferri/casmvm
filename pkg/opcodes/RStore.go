package opcodes

import "strconv"

//RStore load the referenced return stack offset into the evaluation stack
type RStore struct {
	ref uint32
}

func (op *RStore) String() string {
	return "rstore " + strconv.FormatUint(uint64(op.ref), 10)
}

func (op *RStore) Apply(vm VM) VMError {
	val := vm.EvalStack().Pop()
	vm.RetStack().Store(op.ref, val)
	return nil
}

//MakeRStore make an opcode of reference store in the return stack
func MakeRStore(ref uint32) Opcode {
	return &RStore{ref}
}
