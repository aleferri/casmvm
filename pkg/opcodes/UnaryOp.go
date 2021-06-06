package opcodes

import (
	"fmt"
)

//UnaryOp is a unary operation that reference another local variable
type UnaryOp struct {
	name     string
	local    uint16
	ref      uint16
	shape    Shape
	operator UnaryOperator
}

func (op *UnaryOp) String() string {
	return fmt.Sprintf("%%%d = %v %v %%%d", op.local, op.name, op.shape.name, op.ref)
}

func (op *UnaryOp) Apply(vm VM) VMError {
	a := vm.Frame().Values().Peek(op.ref)
	result, err := op.operator(a)
	if err != nil {
		return vm.WrapError(err)
	}
	vm.Frame().Values().Put(op.local, result)
	return nil
}

//MakeUnaryOp create an unary operation to be applied in the stack
func MakeUnaryOp(local uint16, name string, shape Shape, ref uint16, operator UnaryOperator) Opcode {
	return &UnaryOp{name, local, ref, shape, operator}
}
