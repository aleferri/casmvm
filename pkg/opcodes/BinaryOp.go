package opcodes

import (
	"fmt"
	"strconv"
)

//BinaryOp is a binary operation between two references
type BinaryOp struct {
	name     string
	shape    Shape
	local    uint16
	aRef     uint16
	bRef     uint16
	operator BinaryOperator
}

func (op *BinaryOp) String() string {
	t := strconv.FormatUint(uint64(op.local), 10)
	aRef := strconv.FormatUint(uint64(op.aRef), 10)
	bRef := strconv.FormatUint(uint64(op.bRef), 10)
	return fmt.Sprintf("%%%5v = %8v %8v %%%5v %%%5v", t, op.name, op.shape.name, aRef, bRef)
}

func (op *BinaryOp) Apply(vm VM) VMError {
	a := vm.Frame().Values().Peek(op.aRef)
	b := vm.Frame().Values().Peek(op.bRef)
	result, err := op.operator(a, b)
	if err != nil {
		return vm.WrapError(err)
	}
	vm.Frame().Values().Put(op.local, result)
	return nil
}

//MakeBinaryOp create a binary operation to be applied in the stack
func MakeBinaryOp(local uint16, name string, shape Shape, aRef uint16, bRef uint16, operator BinaryOperator) Opcode {
	return &BinaryOp{name, shape, local, aRef, bRef, operator}
}
