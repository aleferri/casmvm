package opcodes

import (
	"fmt"
)

//Branch is the conditional jump opcode: compare the referenced local against a constant, jump by the relative offset if they are equal
type Branch struct {
	cmpval int64
	ifeq   int32
	cmpref uint16
}

func (op *Branch) Locals() []uint16 {
	return []uint16{}
}

func (op *Branch) References() []uint16 {
	r := []uint16{uint16(op.cmpref)}
	return r
}

func (op *Branch) String() string {
	return fmt.Sprintf("ifeq %%%d %d %d", op.cmpref, op.cmpval, op.ifeq)
}

//Offset of the jump, relative to the next instruction
func (op *Branch) Offset() int32 {
	return op.ifeq
}

//WithOffset return a copy of the branch with the provided offset
func (op *Branch) WithOffset(offset int32) *Branch {
	return &Branch{op.cmpval, offset, op.cmpref}
}

func (op *Branch) Apply(vm VM) VMError {
	top := vm.Frame().Values().Peek(op.cmpref)
	if top == op.cmpval {
		vm.Goto(op.ifeq)
	}
	return nil
}

//MakeBranch opcode
func MakeBranch(cmpval int64, cmpref uint16, offset int32) Opcode {
	return &Branch{cmpval: cmpval, ifeq: offset, cmpref: cmpref}
}
