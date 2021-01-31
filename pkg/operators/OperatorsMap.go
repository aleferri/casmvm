package operators

//BinaryOperatorsSymbols for client application
var BinaryOperatorsSymbols = map[string]func(int64, int64) (int64, error){
	"&&": WrapBinaryOp(EvalLogicalAnd), "||": WrapBinaryOp(EvalLogicalOR), "==": WrapBinaryOp(EvalEqual), "!=": WrapBinaryOp(EvalNotEqual),
	">": WrapBinaryOp(EvalGreaterThan), ">=": WrapBinaryOp(EvalGreaterEqualThan), "<=": WrapBinaryOp(EvalLessEqualThan),
	"<": WrapBinaryOp(EvalLessThan), "<<": WrapBinaryOp(EvalShiftLeft), ">>": WrapBinaryOp(EvalShiftRight), "+": WrapBinaryOp(EvalAdd),
	"-": WrapBinaryOp(EvalSub), "*": WrapBinaryOp(EvalMul), "/": EvalDiv, "%": EvalMod, "&": WrapBinaryOp(EvalAnd),
	"|": WrapBinaryOp(EvalOr), "^": WrapBinaryOp(EvalXor),
}

//BinaryOperatorsNames for casmvm
var BinaryOperatorsNames = map[string]func(int64, int64) (int64, error){
	"andl": WrapBinaryOp(EvalLogicalAnd), "orl": WrapBinaryOp(EvalLogicalOR), "cmpeq": WrapBinaryOp(EvalEqual), "cmpne": WrapBinaryOp(EvalNotEqual),
	"cmpgt": WrapBinaryOp(EvalGreaterThan), "cmpge": WrapBinaryOp(EvalGreaterEqualThan), "cmple": WrapBinaryOp(EvalLessEqualThan),
	"cmplt": WrapBinaryOp(EvalLessThan), "shl": WrapBinaryOp(EvalShiftLeft), "shr": WrapBinaryOp(EvalShiftRight),
	"add": WrapBinaryOp(EvalAdd), "sub": WrapBinaryOp(EvalSub), "mul": WrapBinaryOp(EvalMul), "div": EvalDiv, "mod": EvalMod,
	"and": WrapBinaryOp(EvalAnd), "or": WrapBinaryOp(EvalOr), "xor": WrapBinaryOp(EvalXor),
}

//UnaryOperatorsSymbols for client application
var UnaryOperatorsSymbols = map[string]func(int64) (int64, error){
	"-": WrapUnaryOp(EvalNeg), "!": WrapUnaryOp(EvalLogicalNot), "~": WrapUnaryOp(EvalNot),
}

//UnaryOperatorsNames for casmvm
var UnaryOperatorsNames = map[string]func(int64) (int64, error){
	"neg": WrapUnaryOp(EvalNeg), "notl": WrapUnaryOp(EvalLogicalNot), "not": WrapUnaryOp(EvalNot),
}
