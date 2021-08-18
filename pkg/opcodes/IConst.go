package opcodes

import (
	"fmt"
)

//IConst push integer constant
type IConst struct {
	value  int64
	assign uint16
}

func (op *IConst) Locals() []uint16 {
	return []uint16{op.assign}
}

func (op *IConst) References() []uint16 {
	r := []uint16{}
	return r
}

func (op *IConst) Constant() int64 {
	return op.value
}

func (op *IConst) String() string {
	return fmt.Sprintf("%%%d = const i64 %d", op.assign, op.value)
}

func (op *IConst) Apply(vm VM) VMError {
	vm.Frame().Values().Put(op.assign, op.value)
	return nil
}

//MakeIConst to load the specified constant in the eval stack
func MakeIConst(assign uint16, value int64) Opcode {
	return &IConst{value, assign}
}
