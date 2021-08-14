package opcodes

import (
	"fmt"
)

//Branch is the branching opcode, pop the integer constant, check against the compare then branch to the result
type Branch struct {
	cmpval int64
	ifeq   int32
	cmpref uint16
}

func (op *Branch) String() string {
	return fmt.Sprintf("ifeq %%%d %d %d", op.cmpref, op.cmpval, op.ifeq)
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
