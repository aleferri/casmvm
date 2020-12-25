package opcodes

import "strconv"

//Goto is the branching opcode, pop the integer constant, check against the compare then branch to the result
type Goto struct {
	offset int32
}

func (op *Goto) String() string {
	return "goto " + strconv.FormatInt(int64(op.offset), 10)
}

func (op *Goto) Apply(vm VM) VMError {
	vm.GotoOffset(op.offset)
	return nil
}

//MakeGoto opcode
func MakeGoto(offset int32) Opcode {
	return &Goto{offset}
}
