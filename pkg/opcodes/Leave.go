package opcodes

import "fmt"

//Leave a Frame
type Leave struct {
	refs []uint16
}

func (op *Leave) String() string {
	refs := ""
	for _, e := range op.refs {
		refs += fmt.Sprintf("%%%d ", e)
	}
	return "leave " + refs
}

func (op *Leave) Apply(vm VM) VMError {
	frame := vm.Frame()
	for i, v := range op.refs {
		frame.Returns().Put(uint16(i), frame.Local(v))
	}
	vm.Leave()
	return nil
}

//MakeLeave make an opcode of reference load
func MakeLeave(refs ...uint16) Opcode {
	return &Leave{refs}
}
