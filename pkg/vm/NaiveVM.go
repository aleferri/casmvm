package vm

import (
	"fmt"

	"github.com/aleferri/casmvm/pkg/opcodes"
)

type NaiveVM struct {
	ip   uint32
	es   opcodes.Stack
	rs   opcodes.Stack
	list []opcodes.Opcode
	halt bool
}

func (t *NaiveVM) EvalStack() opcodes.Stack {
	return t.es
}

func (t *NaiveVM) RetStack() opcodes.Stack {
	return t.rs
}

func (t *NaiveVM) Goto(ptr uint32) {
	t.ip = ptr
}

func (t *NaiveVM) GotoOffset(disp int32) {
	t.ip = uint32(int32(t.ip) + disp)
}

func (t *NaiveVM) WrapError(e error) opcodes.VMError {
	return &OpcodeError{e, t.ip}
}

func (t *NaiveVM) Halt() {
	t.halt = true
}

func (t *NaiveVM) Pointer() uint32 {
	return t.ip
}

func (t *NaiveVM) Run(debugMode bool) opcodes.VMError {
	var err opcodes.VMError = nil
	for !t.halt && err == nil && int(t.ip) < len(t.list) {
		op := t.list[int(t.ip)]
		t.ip++
		if debugMode {
			fmt.Println(op.String())
			fmt.Println(t.es)
			fmt.Println(t.rs)
		}
		err = op.Apply(t)
	}
	return err
}

func MakeNaiveVM(listing []opcodes.Opcode) *NaiveVM {
	return &NaiveVM{0, MakeStack(), MakeStack(), listing, false}
}
