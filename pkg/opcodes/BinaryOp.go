package opcodes

//BinaryOp is a binary operation on the top and the next of the stack
type BinaryOp struct {
	name     string
	operator BinaryOperator
}

func (op *BinaryOp) String() string {
	return op.name
}

func (op *BinaryOp) Apply(vm VM) VMError {
	top := vm.EvalStack().Pop()
	next := vm.EvalStack().Pop()
	result, err := op.operator(next, top)
	if err != nil {
		return vm.WrapError(err)
	}
	vm.EvalStack().Push(result)
	return nil
}

//MakeBinaryOp create a binary operation to be applied in the stack
func MakeBinaryOp(name string, operator BinaryOperator) Opcode {
	return &BinaryOp{name, operator}
}
