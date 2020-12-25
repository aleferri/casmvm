package opcodes

import "testing"

type testerr struct {
	embed   error
	address uint32
}

func (t *testerr) Error() string {
	return t.embed.Error()
}

func (t *testerr) OpcodeID() int64 {
	return int64(t.address)
}

type testvm struct {
	ip   uint32
	es   Stack
	rs   Stack
	list []Opcode
	halt bool
}

func (t *testvm) EvalStack() Stack {
	return t.es
}

func (t *testvm) RetStack() Stack {
	return t.rs
}

func (t *testvm) Goto(ptr uint32) {
	t.ip = ptr
}

func (t *testvm) GotoOffset(disp int32) {
	t.ip = uint32(int32(t.ip) + disp)
}

func (t *testvm) WrapError(e error) VMError {
	return &testerr{e, t.ip}
}

func (t *testvm) Halt() {
	t.halt = true
}

func (t *testvm) Pointer() uint32 {
	return t.ip
}

type teststack struct {
	content []int64
}

func (s *teststack) Push(value int64) {
	s.content = append(s.content, value)
}

func (s *teststack) Pop() int64 {
	last := len(s.content) - 1
	val := s.content[last]
	s.content = s.content[:last]
	return val
}

func (s *teststack) Empty() bool {
	return len(s.content) == 0
}

func (s *teststack) Load(offset uint32) int64 {
	last := len(s.content) - 1
	return s.content[last-int(offset)]
}

func (s *teststack) Store(offset uint32, value int64) {
	last := len(s.content) - 1
	s.content[last-int(offset)] = value
}

func makeTestStack() Stack {
	return &teststack{[]int64{}}
}

func TestOpcodes(t *testing.T) {
	vm := testvm{0, makeTestStack(), makeTestStack(), []Opcode{}, false}
	vm.list = append(vm.list, MakeIConst(8))
	vm.list = append(vm.list, MakeRPush())
	vm.list = append(vm.list, MakeRLoad(0))
	vm.list = append(vm.list, MakeIConst(16))
	vm.list = append(vm.list, MakeUnaryOp("-", func(a int64) (int64, error) { return -a, nil }))
	vm.list = append(vm.list, MakeBinaryOp("*", func(a, b int64) (int64, error) { return a * b, nil }))
	vm.list = append(vm.list, MakeReturn(1))
	for i, k := range vm.list {
		vm.ip = uint32(i)
		k.Apply(&vm)
	}
	if vm.es.Pop() != -128 {
		t.Fail()
	}
	if !vm.halt {
		t.Fail()
	}
}
