package opcodes

//RPush load the referenced return stack offset into the evaluation stack
type RPush struct {
}

func (op *RPush) String() string {
	return "rpush"
}

func (op *RPush) Apply(vm VM) VMError {
	vm.RetStack().Push(vm.EvalStack().Pop())
	return nil
}

//MakeRPush make an opcode of push to the return stack
func MakeRPush() Opcode {
	return &RPush{}
}
