package fast

import (
	"fmt"

	"github.com/aleferri/casmvm/pkg/opcodes"
	"github.com/aleferri/casmvm/pkg/vmex"
)

func DeadCodeElimination(c vmex.Callable) vmex.Callable {
	listing := c.Listing()
	size := len(listing)
	used := listing[size-1].References()
	aliveReverse := []opcodes.Opcode{listing[size-1]}
	for i := size - 2; i >= 0; i-- {
		op := listing[i]
		skip := len(op.Locals()) > 0
		for _, local := range op.Locals() {
			for _, ref := range used {
				if ref == local {
					skip = false
					goto out
				}
			}
		}
	out:
		if !skip {
			aliveReverse = append(aliveReverse, op)
			used = append(used, op.References()...)
		}
	}
	fmt.Println()
	fmt.Println("_"+c.Name(), ":")
	alive := []opcodes.Opcode{}
	for i := len(aliveReverse) - 1; i >= 0; i-- {
		alive = append(alive, aliveReverse[i])
		fmt.Println(aliveReverse[i])
	}
	fmt.Println()
	return vmex.MakeCallable(c.Name(), c.Params(), alive)
}
