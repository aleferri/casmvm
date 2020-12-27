package operators

//BinaryOperatorsSymbols for client application
var BinaryOperatorsSymbols = map[string]func(int64, int64) (int64, error){
	"&&": EvalLogicalAnd, "||": EvalLogicalOR, "==": EvalEqual, "!=": EvalNotEqual, ">": EvalGreaterThan, ">=": EvalGreaterEqualThan, "<=": EvalLessEqualThan,
	"<": EvalLessThan, "<<": EvalShiftLeft, ">>": EvalShiftRight, "+": EvalAdd, "-": EvalSub, "*": EvalMul, "/": EvalDiv, "%": EvalMod, "&": EvalAnd,
	"|": EvalOr, "^": EvalXor,
}

//BinaryOperatorsNames for casmvm
var BinaryOperatorsNames = map[string]func(int64, int64) (int64, error){
	"andl": EvalLogicalAnd, "orl": EvalLogicalOR, "cmpeq": EvalEqual, "cmpne": EvalNotEqual, "cmpgt": EvalGreaterThan, "cmpge": EvalGreaterEqualThan,
	"cmple": EvalLessEqualThan, "cmplt": EvalLessThan, "shl": EvalShiftLeft, "shr": EvalShiftRight, "add": EvalAdd, "sub": EvalSub, "mul": EvalMul,
	"div": EvalDiv, "mod": EvalMod, "and": EvalAnd, "or": EvalOr, "xor": EvalXor,
}

//UnaryOperatorsSymbols for client application
var UnaryOperatorsSymbols = map[string]func(int64) (int64, error){
	"-": EvalNeg, "!": EvalLogicalNot, "~": EvalNot,
}

//UnaryOperatorsNames for casmvm
var UnaryOperatorsNames = map[string]func(int64) (int64, error){
	"neg": EvalNeg, "notl": EvalLogicalNot, "not": EvalNot,
}
