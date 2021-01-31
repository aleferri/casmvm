package operators

func WrapUnaryOp(fn func(int64) int64) func(int64) (int64, error) {
	return func(a int64) (int64, error) {
		return fn(a), nil
	}
}

func EvalNeg(a int64) int64 {
	return -a
}

func EvalNot(a int64) int64 {
	return ^a
}

func EvalLogicalNot(a int64) int64 {
	if a == 0 {
		return 1
	}
	return 0
}
