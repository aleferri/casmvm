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

func (op *Local) Locals() []uint16 {
	return []uint16{op.local}
}

func (op *Local) References() []uint16 {
	r := []uint16{op.ref}
	return r
}

//Shape of the copied value
func (op *Local) Shape() Shape {
	return op.shape
}

func (op *Local) String() string {
	return fmt.Sprintf("%%%d = local %s %%%d", op.local, op.shape.name, op.ref)
}

func (op *Local) Apply(vm VM) VMError {
	vm.Frame().Values().Put(op.local, op.shape.Reshape(vm.Frame().Local(op.ref)))
	return nil
}

//MakeAssignment make an opcode that assigns a local to another local
func MakeAssignment(local uint16, ref uint16, shape Shape) Opcode {
	return &Local{local, ref, shape}
}
