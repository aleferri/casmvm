package opcodes

//Opcode for the VM interface
type Opcode interface {
	String() string
	Apply(vm VM) VMError
}
