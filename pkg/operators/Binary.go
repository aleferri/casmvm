package operators

import "errors"

func EvalShiftLeft(a int64, b int64) (int64, error) {
	return a << b, nil
}

func EvalShiftRight(a int64, b int64) (int64, error) {
	return a >> b, nil
}

func EvalAdd(a int64, b int64) (int64, error) {
	return a + b, nil
}

func EvalSub(a int64, b int64) (int64, error) {
	return a - b, nil
}

func EvalOr(a int64, b int64) (int64, error) {
	return a | b, nil
}

func EvalMul(a int64, b int64) (int64, error) {
	return a * b, nil
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

func EvalXor(a int64, b int64) (int64, error) {
	return a ^ b, nil
}

func EvalAnd(a int64, b int64) (int64, error) {
	return a & b, nil
}

func EvalLogicalOR(a int64, b int64) (int64, error) {
	if a != 0 || b != 0 {
		return 1, nil
	}
	return 0, nil
}

func EvalLogicalAnd(a int64, b int64) (int64, error) {
	if a != 0 && b != 0 {
		return 1, nil
	}
	return 0, nil
}

func EvalNotEqual(a int64, b int64) (int64, error) {
	if a != b {
		return 1, nil
	}
	return 0, nil
}

func EvalEqual(a int64, b int64) (int64, error) {
	if a == b {
		return 1, nil
	}
	return 0, nil
}

func EvalGreaterEqualThan(a int64, b int64) (int64, error) {
	if a >= b {
		return 1, nil
	}
	return 0, nil
}

func EvalLessEqualThan(a int64, b int64) (int64, error) {
	if a <= b {
		return 1, nil
	}
	return 0, nil
}

func EvalLessThan(a int64, b int64) (int64, error) {
	if a < b {
		return 1, nil
	}
	return 0, nil
}

func EvalGreaterThan(a int64, b int64) (int64, error) {
	if a > b {
		return 1, nil
	}
	return 0, nil
}
