package fast

import (
	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmex"
)

//Fold constants
func Fold(c vmex.Callable) vmex.Callable {
	listing := []opcodes.Opcode{}
	constants := map[uint16]int64{}
	for _, op := range c.Listing() {
		refs := op.References()
		locals := op.Locals()
		unary, isUnary := op.(*opcodes.UnaryOp)
		if isUnary {
			val, ok := constants[refs[0]]
			if ok {
				result, _ := unary.Operator()(val)
				constants[locals[0]] = result
				c := opcodes.MakeIConst(locals[0], result)
				listing = append(listing, c)
			} else {
				listing = append(listing, op)
			}
			continue
		}
		binary, isBinary := op.(*opcodes.BinaryOp)
		if isBinary {
			left, leftOk := constants[binary.Left().Reference()]
			right, rightOk := constants[binary.Right().Reference()]
			if leftOk && rightOk {
				result, _ := binary.Operator()(left, right)
				constants[locals[0]] = result
				c := opcodes.MakeIConst(locals[0], result)
				listing = append(listing, c)
			} else if leftOk {
				binary.SetLeft(opcodes.MakeConstant(left))
				listing = append(listing, binary)
			} else if rightOk {
				binary.SetRight(opcodes.MakeConstant(right))
				listing = append(listing, binary)
			} else {
				listing = append(listing, binary)
			}
			continue
		}
		_, isLocal := op.(*opcodes.Local)
		if isLocal {
			val, ok := constants[refs[0]]
			if ok {
				constants[locals[0]] = val
				c := opcodes.MakeIConst(locals[0], val)
				listing = append(listing, c)
				continue
			}
		}
		c, isConst := op.(*opcodes.IConst)
		if isConst {
			constants[locals[0]] = c.Constant()
		}
		listing = append(listing, op)
	}
	return vmex.MakeCallable(c.Name(), c.Params(), listing)
}
