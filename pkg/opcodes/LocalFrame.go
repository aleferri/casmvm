package opcodes

//Register file that contains local variables, but not results
type RegisterFile struct {
	vals []int64
}

//Put a value in the specified register
func (i *RegisterFile) Put(indx uint16, val int64) {
	for len(i.vals) <= int(indx) {
		i.vals = append(i.vals, 0)
	}
	i.vals[indx] = val
}

//Peek value from the specified registers
func (i *RegisterFile) Peek(indx uint16) int64 {
	return i.vals[indx]
}

func (i *RegisterFile) IsEmpty() bool {
	return len(i.vals) == 0
}

func (i *RegisterFile) Size() int {
	return len(i.vals)
}

//Make a RegisterFile for the LocalFrame
func MakeRegisterFile() RegisterFile {
	return RegisterFile{[]int64{}}
}

//LocalFrame is the computation frame, each operation set a value in his own address
type LocalFrame interface {
	//Value of computation performed at address a
	Values() *RegisterFile
	//Local value of a local parameter
	Local(a uint16) int64
	//Return value
	Returns() *RegisterFile
	//PC is the current program pointer
	PC() uint32
}
