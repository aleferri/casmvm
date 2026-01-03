package vmex

import (
	"testing"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmio"
)

func TestOpcodes(t *testing.T) {
	callable := MakeCallable("test", []string{}, []opcodes.Opcode{})
	callable.list = append(callable.list, opcodes.MakeIConst(0, 8))
	callable.list = append(callable.list, opcodes.MakeAssignment(1, 0, opcodes.ShortShape))
	callable.list = append(callable.list, opcodes.MakeIConst(2, 16))
	callable.list = append(callable.list, opcodes.MakeUnaryOp(3, "-", opcodes.IntShape, 2, func(a int64) (int64, error) { return -a, nil }))
	callable.list = append(callable.list, opcodes.MakeBinaryOp(4, "*", opcodes.IntShape, 3, 1, func(a, b int64) (int64, error) { return a * b, nil }))
	callable.list = append(callable.list, opcodes.MakeLeave(1, 4))

	initialFrame := MakeVMFrame()

	vm := MakeInterpreter([]Callable{callable}, vmio.MakeVMLoggerConsole(vmio.ALL), initialFrame)
	for i, k := range callable.list {
		initialFrame.pc = uint32(i)
		k.Apply(vm)
	}
	if initialFrame.Returns().Peek(1) != -128 {
		t.Errorf("Registers %v\n", initialFrame.values)
		t.Errorf("Unexpected %v\n", initialFrame.Returns())
		t.Fail()
	}
	if !vm.halt {
		t.Errorf("Expected halt\n")
		t.Fail()
	}
}
