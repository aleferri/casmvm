package opcodes

//UnaryOperator is the interface for all unary operators
type UnaryOperator func(a int64) (int64, error)
