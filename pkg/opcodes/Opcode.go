package opcodes

//Opcode for the VM interface
type Opcode interface {
	String() string
	Locals() []uint16
	References() []uint16
	Apply(vm VM) VMError
}
