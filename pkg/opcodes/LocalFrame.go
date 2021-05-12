package opcodes

//Register file that contains local variables, but not results
type RegisterFile struct {
	vals []int64
}

//Put a value in the specified register
func (i *RegisterFile) Put(indx uint16, val int64) {
	cmp := uint16(len(i.vals))
	if cmp < indx {
		old := i.vals
		i.vals = make([]int64, indx+1)
		for _, n := range old {
			i.vals[n] = old[n]
		}
	} else if cmp == indx {
		i.vals = append(i.vals, val)
	} else {
		i.vals[indx] = val
	}
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
	PC() uint16
}
