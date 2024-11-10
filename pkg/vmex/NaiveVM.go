package vmex

import (
	"fmt"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmio"
)

// Callable is a sequential list of opcodes to execute
type Callable struct {
	name   string
	params []string
	list   []opcodes.Opcode
}

// Listing of opcodes
func (c Callable) Listing() []opcodes.Opcode {
	return c.list
}

// Params for the callable
func (c Callable) Params() []string {
	return c.params
}

// Name of the callable
func (c Callable) Name() string {
	return c.name
}

// Dump content of the callable
func (c Callable) Dump() {
	fmt.Println("_"+c.name, ":")
	fmt.Printf("Total: %d opcodes\n", len(c.list))
	for _, op := range c.list {
		fmt.Println(op.String())
	}
}

// Make a callable object for the VM
func MakeCallable(name string, params []string, list []opcodes.Opcode) Callable {
	return Callable{name, params, list}
}

// NaiveVM is a simple implementation of the VM interface found on opcodes
type NaiveVM struct {
	callables []Callable
	logger    vmio.VMLogger
	current   VMFrame
	halt      bool
	leave     bool
	verbose   bool
}

func (t *NaiveVM) Frame() opcodes.LocalFrame {
	return &t.current
}

func (t *NaiveVM) Enter(callIndex int32, vals ...uint16) (opcodes.LocalFrame, opcodes.VMError) {
	return t.Invoke(callIndex, vals...)
}

func (t *NaiveVM) Invoke(callIndex int32, vals ...uint16) (opcodes.LocalFrame, opcodes.VMError) {
	callable := t.callables[callIndex]
	if t.verbose {
		fmt.Printf("Enter callable %d(%s) with args %v\n", callIndex, callable.name, vals)
	}

	prev := t.current
	next := MakeVMFrame()
	for i, v := range vals {
		next.values.Put(uint16(i), prev.Local(v))
	}

	if t.verbose {
		fmt.Println("Accept", next.values)
	}

	t.current = next
	err := t.Run(callable, t.verbose)
	t.current = prev
	t.leave = false

	if t.verbose {
		fmt.Printf("Leave callable %d(%s) with returns %v\n", callIndex, callable.name, next.returns)
	}

	return &next, err
}

func (t *NaiveVM) Leave() {
	t.leave = true
}

func (t *NaiveVM) Goto(disp int32) {
	t.current.pc = uint32(int32(t.current.pc) + disp)
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
	for !t.halt && !t.leave && err == nil && int(t.current.pc) < len(c.list) {
		op := c.list[int(t.current.pc)]
		t.current.pc++
		if debugMode {
			fmt.Println(op.String())
		}
		err = op.Apply(t)
	}
	return err
}

func (t *NaiveVM) Start(fIndex int32, frame opcodes.LocalFrame) opcodes.VMError {
	t.halt = false
	callable := t.callables[fIndex]

	prev := t.current

	t.current = MakeVMFrame()
	t.current.pc = frame.PC()
	t.current.returns = frame.Returns()
	t.current.values = frame.Values()

	err := t.Run(callable, t.verbose)

	t.current = prev
	t.leave = false

	return err
}

func (t *NaiveVM) Dump(index int32) {
	t.callables[index].Dump()
}

func (t *NaiveVM) DumpAll() {
	for index, c := range t.callables {
		fmt.Printf("Callable %d\n", index)
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
	return &NaiveVM{callables, log, bootstrap, false, false, false}
}

func MakeVerboseNaiveVM(callables []Callable, log vmio.VMLogger, bootstrap VMFrame) *NaiveVM {
	return &NaiveVM{callables, log, bootstrap, false, false, true}
}
