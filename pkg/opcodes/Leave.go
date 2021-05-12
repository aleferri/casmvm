package opcodes

import (
	"strconv"
)

type Leave struct {
	refs []uint16
}

func (op *Leave) String() string {
	refs := ""
	for _, e := range op.refs {
		refs += strconv.FormatUint(uint64(e), 10) + " "
	}
	return "        leave  " + refs
}

func (op *Leave) Apply(vm VM) VMError {
	vm.Leave(op.refs...)
	return nil
}

//MakeLeave make an opcode of reference load
func MakeLeave(refs ...uint16) Opcode {
	return &Leave{refs}
}
