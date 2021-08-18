package opcodes

import (
	"strconv"

	"github.com/aleferri/casmvm/pkg/vmio"
)

//SigWarning issue a warning on a specified parameter/local variable
type SigWarning struct {
	msg string
	ref uint16
}

func (op *SigWarning) Locals() []uint16 {
	return []uint16{}
}

func (op *SigWarning) References() []uint16 {
	r := []uint16{op.ref}
	return r
}

func (op *SigWarning) String() string {
	return "sigwarn " + strconv.FormatUint(uint64(op.ref), 10) + " " + op.msg
}

func (op *SigWarning) Apply(vm VM) VMError {
	val := vm.Frame().Values().Peek(op.ref)
	vm.Logger().Log(vmio.WARNING, op.msg+strconv.FormatInt(val, 10))
	return nil
}

//MakeSigWarning make an opcode of reference check
func MakeSigWarning(msg string, ref uint16) Opcode {
	return &SigWarning{msg, ref}
}
