package opcodes

//VMError is the error thrown by the VM
type VMError interface {
	Error() string
	OpcodeID() int64
}
