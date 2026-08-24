package vmex

import (
	"strings"
	"testing"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/operators"
	"github.com/aleferri/casmvm/pkg/vmio"
)

func makeVM(callables ...Callable) *Interpreter {
	return MakeInterpreter(callables, vmio.MakeVMLoggerConsole(vmio.ALL), MakeVMFrame())
}

func TestOpcodes(t *testing.T) {
	list := []opcodes.Opcode{
		opcodes.MakeIConst(0, 8),
		opcodes.MakeAssignment(1, 0, opcodes.ShortShape),
		opcodes.MakeIConst(2, 16),
		opcodes.MakeUnaryOp(3, "neg", opcodes.IntShape, 2, operators.WrapUnaryOp(operators.EvalNeg)),
		opcodes.MakeBinaryOp(4, "mul", opcodes.IntShape, 3, 1, operators.WrapBinaryOp(operators.EvalMul)),
		opcodes.MakeLeave(1, 4),
	}
	vm := makeVM(MakeCallable("test", []string{}, list))

	frame, err := vm.Enter(0)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if frame.Returns().Size() != 2 {
		t.Fatalf("Expected 2 returns, got %d", frame.Returns().Size())
	}
	if frame.Returns().Peek(0) != 8 {
		t.Errorf("Expected 8 as first return, got %d", frame.Returns().Peek(0))
	}
	if frame.Returns().Peek(1) != -128 {
		t.Errorf("Expected -128 as second return, got %d", frame.Returns().Peek(1))
	}
}

func TestUnaryAndLocalRespectShape(t *testing.T) {
	list := []opcodes.Opcode{
		opcodes.MakeIConst(0, 300),
		opcodes.MakeAssignment(1, 0, opcodes.ByteShape),
		opcodes.MakeIConst(2, 130),
		opcodes.MakeUnaryOp(3, "neg", opcodes.ByteShape, 2, operators.WrapUnaryOp(operators.EvalNeg)),
		opcodes.MakeLeave(1, 3),
	}
	vm := makeVM(MakeCallable("shapes", []string{}, list))

	frame, err := vm.Enter(0)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if frame.Returns().Peek(0) != 44 {
		t.Errorf("Expected local i8 of 300 to be 44, got %d", frame.Returns().Peek(0))
	}
	if frame.Returns().Peek(1) != 126 {
		t.Errorf("Expected neg i8 of 130 to be 126, got %d", frame.Returns().Peek(1))
	}
}

func TestBranchLoop(t *testing.T) {
	//sum the integers from 1 to 5 with a countdown loop
	list := []opcodes.Opcode{
		opcodes.MakeIConst(0, 0), //sum
		opcodes.MakeIConst(1, 5), //counter
		opcodes.MakeIConst(2, 1), //step
		opcodes.MakeBinaryOp(0, "add", opcodes.LongShape, 0, 1, operators.WrapBinaryOp(operators.EvalAdd)),
		opcodes.MakeBinaryOp(1, "sub", opcodes.LongShape, 1, 2, operators.WrapBinaryOp(operators.EvalSub)),
		opcodes.MakeBranch(0, 1, 1),
		opcodes.MakeGoto(-4),
		opcodes.MakeLeave(0),
	}
	vm := makeVM(MakeCallable("loop", []string{}, list))

	frame, err := vm.Enter(0)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if frame.Returns().Peek(0) != 15 {
		t.Errorf("Expected 15, got %d", frame.Returns().Peek(0))
	}
}

func TestEnterPropagatesCalleeError(t *testing.T) {
	main := MakeCallable("main", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 42),
		opcodes.MakeEnter([]uint16{1}, 1, []uint16{0}),
		opcodes.MakeLeave(1),
	})
	boom := MakeCallable("boom", []string{"%0"}, []opcodes.Opcode{
		opcodes.MakeSigError("value is bad: ", 0),
		opcodes.MakeLeave(0),
	})
	vm := makeVM(main, boom)

	_, err := vm.Enter(0)
	if err == nil {
		t.Fatal("Expected the callee error to reach the caller")
	}
	if !strings.Contains(err.Error(), "value is bad: 42") {
		t.Errorf("Expected the callee error, got '%s'", err.Error())
	}
}

func TestEnterArityMismatch(t *testing.T) {
	main := MakeCallable("main", []string{}, []opcodes.Opcode{
		opcodes.MakeEnter([]uint16{0}, 1, []uint16{}),
		opcodes.MakeLeave(0),
	})
	two := MakeCallable("two", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 1),
		opcodes.MakeLeave(0, 0),
	})
	vm := makeVM(main, two)

	_, err := vm.Enter(0)
	if err == nil {
		t.Fatal("Expected an arity mismatch error")
	}
	if !strings.Contains(err.Error(), "expected 1 returns, received 2") {
		t.Errorf("Expected an arity mismatch error, got '%s'", err.Error())
	}
}

func TestOpcodeErrorAddress(t *testing.T) {
	list := []opcodes.Opcode{
		opcodes.MakeIConst(0, 1),
		opcodes.MakeIConst(1, 0),
		opcodes.MakeBinaryOp(2, "div", opcodes.LongShape, 0, 1, operators.EvalDiv),
		opcodes.MakeLeave(2),
	}
	vm := makeVM(MakeCallable("divzero", []string{}, list))

	_, err := vm.Enter(0)
	if err == nil {
		t.Fatal("Expected a division by zero error")
	}
	if err.OpcodeID() != 2 {
		t.Errorf("Expected the error at opcode 2, got %d", err.OpcodeID())
	}
}
