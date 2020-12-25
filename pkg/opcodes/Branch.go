package opcodes

import "strconv"

//Branch is the branching opcode, pop the integer constant, check against the compare then branch to the result
type Branch struct {
	cmpval int64
	ifeq   int32
}

func (op *Branch) String() string {
	return "if= " + strconv.FormatInt(op.cmpval, 10) + ", " + strconv.FormatInt(int64(op.ifeq), 10)
}

func (op *Branch) Apply(vm VM) VMError {
	top := vm.EvalStack().Pop()
	if top == op.cmpval {
		vm.GotoOffset(op.ifeq)
	}
	return nil
}

//MakeBranch opcode
func MakeBranch(cmpval int64, offset int32) Opcode {
	return &Branch{cmpval: cmpval, ifeq: offset}
}
