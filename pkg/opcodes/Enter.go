package opcodes

import (
	"fmt"
	"strings"
)

// Enter a frame
type Enter struct {
	rets     []uint16
	callable uint32
	refs     []uint16
}

func (op *Enter) Locals() []uint16 {
	return op.rets
}

func (op *Enter) References() []uint16 {
	return op.refs
}

func (op *Enter) String() string {
	refs := ""
	for _, e := range op.refs {
		refs += fmt.Sprintf("%%%d ", e)
	}
	rets := ""
	for _, e := range op.rets {
		rets += fmt.Sprintf("%%%d ", e)
	}
	return fmt.Sprintf("[%s] = enter %d %s", strings.Trim(rets, " "), op.callable, refs)
}

func (op *Enter) Apply(vm VM) VMError {
	called, err := vm.Enter(int32(op.callable), op.refs...)
	rets := called.Returns().vals

	if len(rets) != len(op.rets) {
		return vm.WrapError(fmt.Errorf("len of formal returns and effective returns diffs, expected %d returns, received %d instead", len(op.rets), len(rets)))
	}

	for i, r := range rets {
		vm.Frame().Values().Put(op.rets[i], r)
	}
	return err
}

// MakeEnter instruction
func MakeEnter(rets []uint16, frame uint32, refs []uint16) Opcode {
	return &Enter{rets, frame, refs}
}
