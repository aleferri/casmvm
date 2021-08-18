package opcodes

import (
	"fmt"
)

//BinaryOp is a binary operation between two references
type BinaryOp struct {
	name     string
	shape    Shape
	local    uint16
	aRef     Operand
	bRef     Operand
	operator BinaryOperator
}

func (op *BinaryOp) Operator() BinaryOperator {
	return op.operator
}

func (op *BinaryOp) Left() Operand {
	return op.aRef
}

func (op *BinaryOp) Right() Operand {
	return op.bRef
}

func (op *BinaryOp) Locals() []uint16 {
	return []uint16{op.local}
}

func (op *BinaryOp) SetLeft(a Operand) {
	op.aRef = a
}

func (op *BinaryOp) SetRight(b Operand) {
	op.bRef = b
}

func (op *BinaryOp) References() []uint16 {
	r := []uint16{}
	if !op.aRef.IsConst() {
		r = append(r, op.aRef.Reference())
	}

	if !op.bRef.IsConst() {
		r = append(r, op.bRef.Reference())
	}

	return r
}

func (op *BinaryOp) String() string {
	return fmt.Sprintf("%%%d = %v %v %s %s", op.local, op.name, op.shape.name, op.aRef.String(), op.bRef.String())
}

func (op *BinaryOp) Apply(vm VM) VMError {
	a := op.aRef.Value(vm)
	b := op.bRef.Value(vm)
	result, err := op.operator(a, b)
	if err != nil {
		return vm.WrapError(err)
	}
	vm.Frame().Values().Put(op.local, result)
	return nil
}

//MakeBinaryOp create a binary operation to be applied in the stack
func MakeBinaryOp(local uint16, name string, shape Shape, aRef uint16, bRef uint16, operator BinaryOperator) Opcode {
	return &BinaryOp{name, shape, local, MakeReference(aRef), MakeReference(bRef), operator}
}

func MakeBinaryOpConstLeft(local uint16, name string, shape Shape, aRef uint16, c int64, operator BinaryOperator) Opcode {
	return &BinaryOp{name, shape, local, MakeReference(aRef), MakeConstant(c), operator}
}

func MakeBinaryOpConstRight(local uint16, name string, shape Shape, c int64, bRef uint16, operator BinaryOperator) Opcode {
	return &BinaryOp{name, shape, local, MakeConstant(c), MakeReference(bRef), operator}
}
