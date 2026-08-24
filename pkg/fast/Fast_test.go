package fast

import (
	"testing"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/operators"
	"github.com/aleferri/casmvm/pkg/vmex"
	"github.com/aleferri/casmvm/pkg/vmio"
)

func run(t *testing.T, callables ...vmex.Callable) *opcodes.RegisterFile {
	t.Helper()
	vm := vmex.MakeInterpreter(callables, vmio.MakeVMLoggerConsole(vmio.ALL), vmex.MakeVMFrame())
	frame, err := vm.Enter(0)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	return frame.Returns()
}

func add() opcodes.BinaryOperator {
	return operators.WrapBinaryOp(operators.EvalAdd)
}

func sub() opcodes.BinaryOperator {
	return operators.WrapBinaryOp(operators.EvalSub)
}

//il programma di riferimento con un loop: somma degli interi da 1 a 5
func loopCallable(dead bool) vmex.Callable {
	list := []opcodes.Opcode{
		opcodes.MakeIConst(0, 0),
		opcodes.MakeIConst(1, 5),
		opcodes.MakeIConst(2, 1),
		opcodes.MakeBinaryOp(0, "add", opcodes.LongShape, 0, 1, add()),
	}
	offset := int32(-4)
	if dead {
		//opcode morto dentro il corpo del loop: DCE deve accorciare il salto all'indietro
		list = append(list, opcodes.MakeIConst(9, 999))
		offset = -5
	}
	list = append(list,
		opcodes.MakeBinaryOp(1, "sub", opcodes.LongShape, 1, 2, sub()),
		opcodes.MakeBranch(0, 1, 1),
		opcodes.MakeGoto(offset),
		opcodes.MakeLeave(0),
	)
	return vmex.MakeCallable("loop", []string{}, list)
}

func TestFoldStraightLine(t *testing.T) {
	src := vmex.MakeCallable("f", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 2),
		opcodes.MakeIConst(1, 3),
		opcodes.MakeBinaryOp(2, "add", opcodes.LongShape, 0, 1, add()),
		opcodes.MakeAssignment(3, 2, opcodes.LongShape),
		opcodes.MakeUnaryOp(4, "neg", opcodes.LongShape, 3, operators.WrapUnaryOp(operators.EvalNeg)),
		opcodes.MakeLeave(4),
	})
	folded := Fold(src)

	if len(folded.Listing()) != len(src.Listing()) {
		t.Fatalf("Fold must not change the listing size: %d != %d", len(folded.Listing()), len(src.Listing()))
	}
	for i := 2; i <= 4; i++ {
		if _, isConst := folded.Listing()[i].(*opcodes.IConst); !isConst {
			t.Errorf("Opcode %d should be folded to a constant, got '%s'", i, folded.Listing()[i].String())
		}
	}
	if run(t, folded).Peek(0) != -5 {
		t.Errorf("Expected -5, got %d", run(t, folded).Peek(0))
	}
}

func TestFoldAppliesShapes(t *testing.T) {
	src := vmex.MakeCallable("f", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 200),
		opcodes.MakeIConst(1, 100),
		opcodes.MakeBinaryOp(2, "add", opcodes.ByteShape, 0, 1, add()),
		opcodes.MakeLeave(2),
	})
	folded := Fold(src)

	original := run(t, src).Peek(0)
	optimized := run(t, folded).Peek(0)
	if original != optimized {
		t.Errorf("Folding changed the result: %d != %d", original, optimized)
	}
	if optimized != 44 {
		t.Errorf("Expected add i8 200 100 to be 44, got %d", optimized)
	}
}

func TestFoldInvalidatesReassignedLocals(t *testing.T) {
	//il chiamato riceve %0 dal chiamante: dopo '%1 = add %0 %0' la costante in %1 non vale più
	callee := vmex.MakeCallable("callee", []string{"%0"}, []opcodes.Opcode{
		opcodes.MakeIConst(1, 7),
		opcodes.MakeBinaryOp(1, "add", opcodes.LongShape, 0, 0, add()),
		opcodes.MakeBinaryOp(2, "add", opcodes.LongShape, 1, 1, add()),
		opcodes.MakeLeave(2),
	})
	folded := Fold(callee)

	if _, isConst := folded.Listing()[2].(*opcodes.IConst); isConst {
		t.Fatal("Opcode 2 depends on a reassigned local, it must not be folded")
	}

	main := vmex.MakeCallable("main", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 3),
		opcodes.MakeEnter([]uint16{1}, 1, []uint16{0}),
		opcodes.MakeLeave(1),
	})
	if got := run(t, main, folded).Peek(0); got != 12 {
		t.Errorf("Expected (3+3)+(3+3) = 12, got %d", got)
	}
}

func TestFoldClearsFactsAtJumpTargets(t *testing.T) {
	src := loopCallable(false)
	folded := Fold(src)

	if len(folded.Listing()) != len(src.Listing()) {
		t.Fatalf("Fold must not change the listing size")
	}
	//l'add in testa al loop è un target di salto: i fatti vanno scartati, l'opcode resta
	if _, isConst := folded.Listing()[3].(*opcodes.IConst); isConst {
		t.Fatal("The loop body must not be folded, its locals change at every iteration")
	}
	if got := run(t, folded).Peek(0); got != 15 {
		t.Errorf("Expected 15, got %d", got)
	}
}

func TestFoldKeepsRuntimeErrors(t *testing.T) {
	src := vmex.MakeCallable("f", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 1),
		opcodes.MakeIConst(1, 0),
		opcodes.MakeBinaryOp(2, "div", opcodes.LongShape, 0, 1, operators.EvalDiv),
		opcodes.MakeLeave(2),
	})
	folded := Fold(src)

	if _, isConst := folded.Listing()[2].(*opcodes.IConst); isConst {
		t.Fatal("A failing operation must not be folded, the error belongs to the runtime")
	}
	vm := vmex.MakeInterpreter([]vmex.Callable{folded}, vmio.MakeVMLoggerConsole(vmio.ALL), vmex.MakeVMFrame())
	if _, err := vm.Enter(0); err == nil {
		t.Fatal("Expected the division by zero to survive the folding")
	}
}

func TestFoldDoesNotMutateTheInput(t *testing.T) {
	binary := opcodes.MakeBinaryOp(2, "add", opcodes.LongShape, 1, 0, add()).(*opcodes.BinaryOp)
	callee := vmex.MakeCallable("callee", []string{"%0"}, []opcodes.Opcode{
		opcodes.MakeIConst(1, 2),
		binary,
		opcodes.MakeLeave(2),
	})
	Fold(callee)

	if binary.Left().IsConst() || binary.Right().IsConst() {
		t.Fatal("Fold must not mutate the operands of the input callable")
	}
}

func TestDCERemovesDeadCodeAndKeepsSideEffects(t *testing.T) {
	src := vmex.MakeCallable("f", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 1),
		opcodes.MakeIConst(1, 999), //morto
		opcodes.MakeSigWarning("just saying: ", 0),
		opcodes.MakeIConst(2, 777), //morto
		opcodes.MakeLeave(0),
	})
	pruned := DeadCodeElimination(src)

	if len(pruned.Listing()) != 3 {
		t.Fatalf("Expected 3 opcodes after DCE, got %d", len(pruned.Listing()))
	}
	if got := run(t, pruned).Peek(0); got != 1 {
		t.Errorf("Expected 1, got %d", got)
	}
}

func TestDCERebasesBackwardJumps(t *testing.T) {
	src := loopCallable(true)
	pruned := DeadCodeElimination(src)

	if len(pruned.Listing()) != len(src.Listing())-1 {
		t.Fatalf("Expected the dead opcode to be removed, got %d opcodes", len(pruned.Listing()))
	}
	original := run(t, src).Peek(0)
	optimized := run(t, pruned).Peek(0)
	if original != optimized || optimized != 15 {
		t.Errorf("DCE changed the result: original %d, optimized %d", original, optimized)
	}
}

func TestDCERebasesForwardJumps(t *testing.T) {
	src := vmex.MakeCallable("f", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 1),
		opcodes.MakeBranch(1, 0, 2), //salta i due opcode morti
		opcodes.MakeIConst(1, 999),  //morto
		opcodes.MakeIConst(2, 777),  //morto
		opcodes.MakeLeave(0),
	})
	pruned := DeadCodeElimination(src)

	if len(pruned.Listing()) != 3 {
		t.Fatalf("Expected 3 opcodes after DCE, got %d", len(pruned.Listing()))
	}
	if got := run(t, pruned).Peek(0); got != 1 {
		t.Errorf("Expected 1, got %d", got)
	}
}

func TestDCEKeepsFailingOperations(t *testing.T) {
	src := vmex.MakeCallable("f", []string{}, []opcodes.Opcode{
		opcodes.MakeIConst(0, 1),
		opcodes.MakeIConst(1, 0),
		opcodes.MakeBinaryOp(2, "div", opcodes.LongShape, 0, 1, operators.EvalDiv), //risultato inutilizzato
		opcodes.MakeIConst(3, 7),
		opcodes.MakeLeave(3),
	})
	pruned := DeadCodeElimination(src)

	vm := vmex.MakeInterpreter([]vmex.Callable{pruned}, vmio.MakeVMLoggerConsole(vmio.ALL), vmex.MakeVMFrame())
	if _, err := vm.Enter(0); err == nil {
		t.Fatal("Expected the division by zero to survive DCE even if its result is unused")
	}
}

func TestFoldThenDCEOnLoop(t *testing.T) {
	src := loopCallable(true)
	optimized := DeadCodeElimination(Fold(src))

	original := run(t, src).Peek(0)
	final := run(t, optimized).Peek(0)
	if original != final || final != 15 {
		t.Errorf("The pipeline changed the result: original %d, optimized %d", original, final)
	}
}
