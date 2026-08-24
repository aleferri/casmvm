package fast

import (
	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmex"
)

//hasSideEffects report whether the opcode does anything beside writing its locals
func hasSideEffects(op opcodes.Opcode) bool {
	switch typed := op.(type) {
	case *opcodes.Enter, *opcodes.Leave, *opcodes.Branch, *opcodes.Goto, *opcodes.SigError, *opcodes.SigWarning:
		return true
	case *opcodes.BinaryOp:
		//div and mod can stop the execution with an error, dropping them would drop the error too
		name := typed.Name()
		return name == "div" || name == "mod"
	}
	return false
}

//DeadCodeElimination drop the opcodes whose results are never used.
//Opcodes with side effects are the roots of the liveness analysis and are
//always kept; the offsets of the surviving jumps are rebased to account for
//the removed opcodes. The input callable is not modified.
func DeadCodeElimination(c vmex.Callable) vmex.Callable {
	listing := c.Listing()
	size := len(listing)

	live := make([]bool, size)
	for i, op := range listing {
		live[i] = hasSideEffects(op)
	}

	//fixpoint: an opcode is live when any of its locals is referenced by a live opcode
	used := map[uint16]bool{}
	for changed := true; changed; {
		changed = false
		for i, op := range listing {
			if !live[i] {
				continue
			}
			for _, ref := range op.References() {
				if !used[ref] {
					used[ref] = true
					changed = true
				}
			}
		}
		for i, op := range listing {
			if live[i] {
				continue
			}
			for _, local := range op.Locals() {
				if used[local] {
					live[i] = true
					changed = true
					break
				}
			}
		}
	}

	//prefix[k] is the number of live opcodes before index k; prefix[size] is the new listing size.
	//For any old index k, prefix[k] is also the new index of the first live opcode at or after k.
	prefix := make([]int, size+1)
	for i := 0; i < size; i++ {
		prefix[i+1] = prefix[i]
		if live[i] {
			prefix[i+1]++
		}
	}

	rebase := func(index int, offset int32) int32 {
		target := index + 1 + int(offset)
		if target < 0 || target > size {
			//an out of range jump ends the callable, exactly like jumping past the last opcode
			target = size
		}
		return int32(prefix[target] - (prefix[index] + 1))
	}

	alive := make([]opcodes.Opcode, 0, prefix[size])
	for i, op := range listing {
		if !live[i] {
			continue
		}
		switch jump := op.(type) {
		case *opcodes.Branch:
			alive = append(alive, jump.WithOffset(rebase(i, jump.Offset())))
		case *opcodes.Goto:
			alive = append(alive, jump.WithOffset(rebase(i, jump.Offset())))
		default:
			alive = append(alive, op)
		}
	}
	return vmex.MakeCallable(c.Name(), c.Params(), alive)
}
