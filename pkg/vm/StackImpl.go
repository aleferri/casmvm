package vm

import "github.com/aleferri/casmvm/pkg/opcodes"

//StackImpl is the default stack implementation
type StackImpl struct {
	content []int64
}

func (s *StackImpl) Push(value int64) {
	s.content = append(s.content, value)
}

func (s *StackImpl) Pop() int64 {
	last := len(s.content) - 1
	val := s.content[last]
	s.content = s.content[:last]
	return val
}

func (s *StackImpl) Empty() bool {
	return len(s.content) == 0
}

func (s *StackImpl) Load(offset uint32) int64 {
	last := len(s.content) - 1
	return s.content[last-int(offset)]
}

func (s *StackImpl) Store(offset uint32, value int64) {
	last := len(s.content) - 1
	s.content[last-int(offset)] = value
}

//MakeStack return a stack implementation provided by this package
func MakeStack() opcodes.Stack {
	return &StackImpl{[]int64{}}
}
