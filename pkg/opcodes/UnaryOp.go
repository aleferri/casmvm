package opcodes

//UnaryOp is a binary operation on the top and the next of the stack
type UnaryOp struct {
	name     string
	operator UnaryOperator
}

func (op *UnaryOp) String() string {
	return op.name
}

func (op *UnaryOp) Apply(vm VM) VMError {
	top := vm.EvalStack().Pop()
	result, err := op.operator(top)
	if err != nil {
		return vm.WrapError(err)
	}
	vm.EvalStack().Push(result)
	return nil
}

//MakeUnaryOp create an unary operation to be applied in the stack
func MakeUnaryOp(name string, operator UnaryOperator) Opcode {
	return &UnaryOp{name, operator}
}
