package opcodes

import "strconv"

//IConst push integer constant
type IConst struct {
	value int64
}

func (op *IConst) String() string {
	return "iconst " + strconv.FormatInt(op.value, 10)
}

func (op *IConst) Apply(vm VM) VMError {
	vm.EvalStack().Push(op.value)
	return nil
}

//MakeIConst to load the specified constant in the eval stack
func MakeIConst(value int64) Opcode {
	return &IConst{value}
}
