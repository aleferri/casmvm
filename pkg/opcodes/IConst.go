package opcodes

import (
	"fmt"
	"strconv"
)

//IConst push integer constant
type IConst struct {
	value  int64
	assign uint16
}

func (op *IConst) String() string {
	return fmt.Sprintf("%%%v = const i64 %d", strconv.FormatUint(uint64(op.assign), 10), op.value)
}

func (op *IConst) Apply(vm VM) VMError {
	vm.Frame().Values().Put(op.assign, op.value)
	return nil
}

//MakeIConst to load the specified constant in the eval stack
func MakeIConst(assign uint16, value int64) Opcode {
	return &IConst{value, assign}
}
