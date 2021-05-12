package opcodes

import (
	"strconv"
)

type Enter struct {
	start uint16
	frame uint32
	refs  []uint16
}

func (op *Enter) String() string {
	refs := ""
	for _, e := range op.refs {
		refs += strconv.FormatUint(uint64(e), 10) + " "
	}
	return "        enter " + strconv.FormatUint(uint64(op.start), 10) + strconv.FormatUint(uint64(op.frame), 10) + " " + refs
}

func (op *Enter) Apply(vm VM) VMError {
	called, err := vm.Enter(int32(op.frame), op.refs...)
	rets := called.Returns().vals
	for i, r := range rets {
		vm.Frame().Values().Put(op.start+uint16(i), r)
	}
	return err
}

//MakeEnter instruction
func MakeEnter(start uint16, frame uint32, refs []uint16) Opcode {
	return &Enter{start, frame, refs}
}
