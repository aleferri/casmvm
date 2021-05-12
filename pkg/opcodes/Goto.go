package opcodes

import (
	"fmt"
)

//Goto is the branching opcode, pop the integer constant, check against the compare then branch to the result
type Goto struct {
	offset int32
}

func (op *Goto) String() string {
	return fmt.Sprintf("        goto %d", op.offset)
}

func (op *Goto) Apply(vm VM) VMError {
	vm.Goto(op.offset)
	return nil
}

//MakeGoto opcode
func MakeGoto(offset int32) Opcode {
	return &Goto{offset}
}
