package opcodes

//BinaryOperator is the interface for all binary operators
type BinaryOperator func(a int64, b int64) (int64, error)
