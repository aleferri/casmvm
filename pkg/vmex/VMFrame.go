package vmex

import "github.com/aleferri/casmvm/pkg/opcodes"

//VMFrame keep the current frame state
type VMFrame struct {
	values  *opcodes.RegisterFile
	returns *opcodes.RegisterFile
	pc      uint16
}

func (f *VMFrame) Values() *opcodes.RegisterFile {
	return f.values
}

func (f *VMFrame) Local(a uint16) int64 {
	return f.values.Peek(a)
}

func (f *VMFrame) Returns() *opcodes.RegisterFile {
	return f.returns
}

func (f *VMFrame) PC() uint16 {
	return f.pc
}

func MakeVMFrame() VMFrame {
	locals := opcodes.MakeRegisterFile()
	returns := opcodes.MakeRegisterFile()
	return VMFrame{&locals, &returns, 0}
}
