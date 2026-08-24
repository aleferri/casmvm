package opcodes

import (
	"fmt"
)

//Goto is the unconditional jump opcode: jump by the relative offset
type Goto struct {
	offset int32
}

func (op *Goto) Locals() []uint16 {
	return []uint16{}
}

func (op *Goto) References() []uint16 {
	r := []uint16{}
	return r
}

func (op *Goto) String() string {
	return fmt.Sprintf("goto %d", op.offset)
}

//Offset of the jump, relative to the next instruction
func (op *Goto) Offset() int32 {
	return op.offset
}

//WithOffset return a copy of the goto with the provided offset
func (op *Goto) WithOffset(offset int32) *Goto {
	return &Goto{offset}
}

func (op *Goto) Apply(vm VM) VMError {
	vm.Goto(op.offset)
	return nil
}

//MakeGoto opcode
func MakeGoto(offset int32) Opcode {
	return &Goto{offset}
}
