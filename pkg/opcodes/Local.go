package opcodes

import (
	"fmt"
)

//Local assign a new value to a local variable
type Local struct {
	local uint16
	ref   uint16
	shape Shape
}

func (op *Local) String() string {
	return fmt.Sprintf("%%%d = local i16 %d", op.local, op.ref)
}

func (op *Local) Apply(vm VM) VMError {
	vm.Frame().Values().Put(op.local, vm.Frame().Local(op.ref))
	return nil
}

//MakeRLoad make an opcode of reference load
func MakeAssignment(local uint16, ref uint16, shape Shape) Opcode {
	return &Local{local, ref, shape}
}
