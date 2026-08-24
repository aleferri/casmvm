package fast

import (
	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmex"
)

//jumpTargets collect the indices reachable by a jump in the listing
func jumpTargets(listing []opcodes.Opcode) map[int]bool {
	targets := map[int]bool{}
	for i, op := range listing {
		switch jump := op.(type) {
		case *opcodes.Branch:
			targets[i+1+int(jump.Offset())] = true
		case *opcodes.Goto:
			targets[i+1+int(jump.Offset())] = true
		}
	}
	return targets
}

//operandConstant resolve the operand either as an inline constant or as a local known to hold one
func operandConstant(o opcodes.Operand, constants map[uint16]int64) (int64, bool) {
	if o.IsConst() {
		return o.Value(nil), true
	}
	val, known := constants[o.Reference()]
	return val, known
}

//invalidate forget the constants held by the locals overwritten by the opcode
func invalidate(op opcodes.Opcode, constants map[uint16]int64) {
	for _, local := range op.Locals() {
		delete(constants, local)
	}
}

//Fold replace locals known to hold constants with the constants themselves.
//The pass never adds nor removes opcodes, so jump offsets stay valid; known
//facts are dropped at every jump target and whenever a local is overwritten
//by a value that is not a compile time constant. Operations whose operator
//fails on the known constants are kept untouched: the error belongs to the
//runtime. The input callable is not modified.
func Fold(c vmex.Callable) vmex.Callable {
	src := c.Listing()
	targets := jumpTargets(src)
	listing := make([]opcodes.Opcode, 0, len(src))
	constants := map[uint16]int64{}
	for i, op := range src {
		if targets[i] {
			//a jump may land here, facts collected so far are no longer guaranteed
			constants = map[uint16]int64{}
		}
		listing = append(listing, foldOpcode(op, constants))
	}
	return vmex.MakeCallable(c.Name(), c.Params(), listing)
}

func foldOpcode(op opcodes.Opcode, constants map[uint16]int64) opcodes.Opcode {
	switch typed := op.(type) {
	case *opcodes.IConst:
		constants[typed.Locals()[0]] = typed.Constant()
		return op
	case *opcodes.Local:
		local := typed.Locals()[0]
		val, known := constants[typed.References()[0]]
		if !known {
			delete(constants, local)
			return op
		}
		folded := typed.Shape().Reshape(val)
		constants[local] = folded
		return opcodes.MakeIConst(local, folded)
	case *opcodes.UnaryOp:
		local := typed.Locals()[0]
		val, known := constants[typed.References()[0]]
		if !known {
			delete(constants, local)
			return op
		}
		result, err := typed.Operator()(val)
		if err != nil {
			delete(constants, local)
			return op
		}
		folded := typed.Shape().Reshape(result)
		constants[local] = folded
		return opcodes.MakeIConst(local, folded)
	case *opcodes.BinaryOp:
		return foldBinaryOp(typed, constants)
	default:
		invalidate(op, constants)
		return op
	}
}

func foldBinaryOp(op *opcodes.BinaryOp, constants map[uint16]int64) opcodes.Opcode {
	local := op.Locals()[0]
	left, leftKnown := operandConstant(op.Left(), constants)
	right, rightKnown := operandConstant(op.Right(), constants)
	if leftKnown && rightKnown {
		result, err := op.Operator()(left, right)
		if err == nil {
			folded := op.Shape().Reshape(result)
			constants[local] = folded
			return opcodes.MakeIConst(local, folded)
		}
	}
	delete(constants, local)
	if leftKnown && !op.Left().IsConst() {
		op = op.WithOperands(opcodes.MakeConstant(left), op.Right())
	}
	if rightKnown && !op.Right().IsConst() {
		op = op.WithOperands(op.Left(), opcodes.MakeConstant(right))
	}
	return op
}
