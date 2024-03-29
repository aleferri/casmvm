package vmex

//OpcodeError is an error occurred during opcode execution
type OpcodeError struct {
	embed   error
	address uint32
}

func (t *OpcodeError) Error() string {
	return t.embed.Error()
}

func (t *OpcodeError) OpcodeID() int64 {
	return int64(t.address)
}
