package opcodes

import (
	"fmt"
	"math"
)

//Operand value const or reference
type Operand interface {
	Value(vm VM) int64
	Reference() uint16
	IsConst() bool
	String() string
}

type Constant struct {
	val int64
}

func (c *Constant) String() string {
	return fmt.Sprintf("%d", c.val)
}

func (c *Constant) Value(vm VM) int64 {
	return c.val
}

func (c *Constant) Reference() uint16 {
	return math.MaxUint16
}

func (c *Constant) IsConst() bool {
	return true
}

func MakeConstant(a int64) *Constant {
	return &Constant{a}
}

type Reference struct {
	ref uint16
}

func (r *Reference) String() string {
	return fmt.Sprintf("%%%d", r.ref)
}

func (r *Reference) Value(vm VM) int64 {
	return vm.Frame().Local(r.ref)
}

func (r *Reference) IsConst() bool {
	return false
}

func (r *Reference) Reference() uint16 {
	return r.ref
}

func MakeReference(a uint16) *Reference {
	return &Reference{a}
}
