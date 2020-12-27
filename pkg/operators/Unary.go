package operators

func EvalNeg(a int64) (int64, error) {
	return -a, nil
}

func EvalNot(a int64) (int64, error) {
	return ^a, nil
}

func EvalLogicalNot(a int64) (int64, error) {
	if a == 0 {
		return 1, nil
	}
	return 0, nil
}
