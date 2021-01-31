package operators

import "errors"

//WrapBinaryOp that cannot fail
func WrapBinaryOp(fn func(int64, int64) int64) func(int64, int64) (int64, error) {
	return func(a, b int64) (int64, error) {
		return fn(a, b), nil
	}
}

//EvalShiftLeft on int a and b
func EvalShiftLeft(a int64, b int64) int64 {
	return a << b
}

//EvalShiftRight on int a and b
func EvalShiftRight(a int64, b int64) int64 {
	return a >> b
}

//EvalShiftAdd on int a and b
func EvalAdd(a int64, b int64) int64 {
	return a + b
}

func EvalSub(a int64, b int64) int64 {
	return a - b
}

func EvalOr(a int64, b int64) int64 {
	return a | b
}

func EvalMul(a int64, b int64) int64 {
	return a * b
}

func EvalDiv(a int64, b int64) (int64, error) {
	if b == 0 {
		return 0, errors.New("Division by zero")
	}
	return a / b, nil
}

func EvalMod(a int64, b int64) (int64, error) {
	if b == 0 {
		return 0, errors.New("Division by zero")
	}
	return a % b, nil
}

func EvalXor(a int64, b int64) int64 {
	return a ^ b
}

func EvalAnd(a int64, b int64) int64 {
	return a & b
}

func EvalLogicalOR(a int64, b int64) int64 {
	if a != 0 || b != 0 {
		return 1
	}
	return 0
}

func EvalLogicalAnd(a int64, b int64) int64 {
	if a != 0 && b != 0 {
		return 1
	}
	return 0
}

func EvalNotEqual(a int64, b int64) int64 {
	if a != b {
		return 1
	}
	return 0
}

func EvalEqual(a int64, b int64) int64 {
	if a == b {
		return 1
	}
	return 0
}

func EvalGreaterEqualThan(a int64, b int64) int64 {
	if a >= b {
		return 1
	}
	return 0
}

func EvalLessEqualThan(a int64, b int64) int64 {
	if a <= b {
		return 1
	}
	return 0
}

func EvalLessThan(a int64, b int64) int64 {
	if a < b {
		return 1
	}
	return 0
}

func EvalGreaterThan(a int64, b int64) int64 {
	if a > b {
		return 1
	}
	return 0
}
