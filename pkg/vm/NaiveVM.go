package vm

import (
	"fmt"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmio"
)

//Callable is a sequential list of opcodes to execute
type Callable struct {
	list []opcodes.Opcode
}

func (c *Callable) Set(list []opcodes.Opcode) {
	c.list = list
}

func (c *Callable) Get() []opcodes.Opcode {
	return c.list
}

func (c *Callable) Dump() {
	fmt.Printf("Total: %d opcodes\n", len(c.list))
	for _, op := range c.list {
		fmt.Println(op.String())
	}
}

//Make a callable object for the VM
func MakeCallable(list []opcodes.Opcode) Callable {
	return Callable{list}
}

//NaiveVM is a simple implementation of the VM interface found on opcodes
type NaiveVM struct {
	callables []Callable
	logger    vmio.VMLogger
	current   VMFrame
	halt      bool
	last      bool
	verbose   bool
}

func (t *NaiveVM) Frame() opcodes.LocalFrame {
	return &t.current
}

func (t *NaiveVM) Enter(frame int32, vals ...uint16) (opcodes.LocalFrame, opcodes.VMError) {
	prev := t.current
	wasLast := t.last
	next := MakeVMFrame()
	for i, v := range vals {
		next.values.Put(uint16(i), prev.Local(v))
	}
	list := t.callables[frame]
	err := t.Run(list, t.verbose)
	t.current = prev
	t.last = wasLast
	return &next, err
}

func (t *NaiveVM) Leave(vals ...uint16) {
	for i, v := range vals {
		t.current.returns.Put(uint16(i), t.current.Local(v))
	}
	t.halt = true
}

func (t *NaiveVM) Goto(disp int32) {
	t.current.pc = uint16(int32(t.current.pc) + disp)
}

func (t *NaiveVM) WrapError(e error) opcodes.VMError {
	return &OpcodeError{e, uint32(t.current.pc)}
}

func (t *NaiveVM) Halt() {
	t.halt = true
}

func (t *NaiveVM) Pointer() uint32 {
	return uint32(t.current.PC())
}

func (t *NaiveVM) Run(c Callable, debugMode bool) opcodes.VMError {
	var err opcodes.VMError = nil
	for !t.halt && err == nil && int(t.current.pc) < len(c.list) {
		op := c.list[int(t.current.pc)]
		t.current.pc++
		if debugMode {
			fmt.Println(op.String())
		}
		err = op.Apply(t)
	}
	return err
}

func (t *NaiveVM) Dump(frame int32) {
	t.callables[frame].Dump()
}

func (t *NaiveVM) DumpAll() {
	for frame, c := range t.callables {
		fmt.Printf("Frame %d\n", frame)
		c.Dump()
	}
}

func (t *NaiveVM) Logger() vmio.VMLogger {
	return t.logger
}

func (t *NaiveVM) Callables() []Callable {
	return t.callables
}

func MakeNaiveVM(callables []Callable, log vmio.VMLogger, bootstrap VMFrame) *NaiveVM {
	return &NaiveVM{callables, log, bootstrap, false, true, false}
}

func MakeVerboseNaiveVM(callables []Callable, log vmio.VMLogger, bootstrap VMFrame) *NaiveVM {
	return &NaiveVM{callables, log, bootstrap, false, true, true}
}
