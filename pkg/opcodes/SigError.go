package opcodes

import (
	"errors"
	"strconv"

	"github.com/aleferri/casmvm/pkg/vmio"
)

//SigError issue a warning on a specified parameter/local variable
type SigError struct {
	msg string
	ref uint32
}

func (op *SigError) String() string {
	return "sigerr " + strconv.FormatUint(uint64(op.ref), 10)
}

func (op *SigError) Apply(vm VM) VMError {
	val := vm.RetStack().Load(op.ref)
	vm.Logger().Log(vmio.ERROR, op.msg+strconv.FormatInt(val, 10))
	return vm.WrapError(errors.New(op.msg + strconv.FormatInt(val, 10)))
}

//MakeSigError make an opcode of reference check
func MakeSigError(msg string, ref uint32) Opcode {
	return &SigError{msg, ref}
}
