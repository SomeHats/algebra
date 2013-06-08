package algebra

import (
	"errors"
)

/*

Standard operations:
+|: f(x) +|- g(x)   ->   f'(x) +|- g'(x)
*: f(x) * g(x)      ->   f'(x)*g(x) + f(x)*g'(x)
/: f(x) / g(x)      ->   (f'(x)*g(x) - f(x)*g'(x)) / (g(x)^2)
^: f(x) ^ g(x)      ->   (f(x)^(g(x)-1))*(g(x)*f'(x)+f(x)*ln(f(x))*g'(x))

*/

func (exp *Expression) Differentiate(respect string) (*Expression, error) {
	opType := exp.Type
	op := exp.Op

	switch opType {
	case NUMBER, CONSTANT:
		return &Expression{"0", NUMBER, nil, nil}, nil

	case VARIABLE:
		if op == respect {
			return &Expression{"1", NUMBER, nil, nil}, nil
		} else {
			return &Expression{"0", NUMBER, nil, nil}, nil
		}

	case OP_LOW, OP_MED, OP_HIGH:
		f := exp.Left
		g := exp.Right
		fdash, err := exp.Left.Differentiate(respect)
		gdash, err2 := exp.Right.Differentiate(respect)

		if err != nil || err2 != nil {
			if err == nil {
				return nil, err2
			} else {
				return nil, err
			}
		}

		switch op {
		case "+", "-":
			return  &Expression{op, OP_LOW, fdash, gdash}, nil

		case "*":
			return 	add(
								mul(fdash, g),
								mul(f, gdash)), nil

		case "/":
			return 	div(
								sub(mul(fdash, g), mul(f, gdash)),
								powen(g, "2")), nil

		case "^":
			return 	mul(
								pow(f, suben(g, "1")),
								add(
									mul(g, fdash),
									mul(f, mul(gdash, apply("ln", f))))), nil

		// end of OP_(LOW|MED|HIGH) bits
		}

	case FUNC_PREFIX:
		f := exp.Left
		fdash, err := exp.Left.Differentiate(respect)

		if err != nil {
			return nil, err
		}

		switch op {

		case "sin":
			// sin(f(x))		->	f'(x)*cos(f(x))
			return	mul(
								fdash,
								apply("cos", f)), nil

		case "cos":
			// cos(f(x))		->	-f'(x)*sin(f(x))
			return	mul(
								neg(fdash),
								apply("sin", f)), nil

		case "tan":
			// tan(f(x))		->	f'(x)*sec^2(f(x))
			return 	mul(
								fdash,
								powen(
									apply("sec", f), "2")), nil

		case "sec":
			// sec(f(x))		->	f'(x)*tan(f(x))sec(f(x))
			return 	mul(
								fdash,
								mul(
									apply("tan", f),
									apply("sec", f))), nil

		case "csc", "cosec":
			// cosec(f(x))	->	-f'(x)*cot(f(x))cosec(f(x))
			return 	mul(
								neg(fdash),
								mul(
									apply("cot", f),
									apply("cosec", f))), nil

		case "cot":
			// cot(f(x))		->	-f'(x)*cosec^2(f(x))
			return	mul(
								neg(fdash),
								powen(
									apply("cosec", f),
									"2")), nil

		case "ln":
			// ln(f(x))			->	f'(x)/f(x)
			return	div(fdash, f), nil

		case "log":
			// log(10, f(x))->	f'(x)/(f(x)*ln10)
			return	div(
								fdash,
								mul(
									apply("ln", no("10")),
									f)), nil

		case "sqrt":
			// sqrt(f(x))		->	f'(x)/(2*sqrt(f(x)))
			return	div(
								fdash,
								muln(
									"2",
									apply("sqrt", f))), nil

		case "asin", "arsin", "arcsin":
			// asin(f(x))		->	f'(x)/sqrt(1-f(x)^2)
			return	div(
								fdash,
								apply("sqrt",
									subne(
										"1",
										powen(
											f,
											"2")))), nil

		case "acos", "arcos", "arccos":
			// acos(f(x))		->	-f'(x)/sqrt(1-f(x)^2)
			return	div(
								neg(fdash),
								apply("sqrt",
									subne(
										"1",
										powen(
											f,
											"2")))), nil

		case "atan", "artan", "arctan":
			// atan(f(x))		->	f'(x)/sqrt(f(x)^2+1)
			return	div(
								fdash,
								addn(
									"1",
									powen(
										f,
										"2"))), nil

		case "asec", "arsec", "arcsec":
			return	div(
								fdash,
								mul(
									apply("sqrt",
										subne(
											"1",
											divne(
												"1",
												powen(
													f,
													"2")))),
									powen(
										f,
										"2"))), nil

		case "acsc", "arcsc", "arccsc", "acosec", "arcosec", "arccosec":
			return	div(
								neg(fdash),
								mul(
									apply("sqrt",
										subne(
											"1",
											divne(
												"1",
												powen(
													f,
													"2")))),
									powen(
										f,
										"2"))), nil

			case "acot", "arcot", "arccot":
				return	div(
									neg(fdash),
									addn(
										"1",
										powen(
											f,
											"2"))), nil

			case "sinh":
				return	mul(
									fdash,
									apply("cosh", f)), nil

			case "cosh":
				return	mul(
									fdash,
									apply("sinh", f)), nil

			case "tanh":
				return	mul(
									fdash,
									powen(
										apply("sech", f),
										"2")), nil

			case "sech":
				return	mul(
									neg(fdash),
									mul(
										apply("tanh", f),
										apply("sech", f))), nil

			case "csch", "cosech":
				return	mul(
									neg(fdash),
									mul(
										apply("coth", f),
										apply("cosech", f))), nil

			case "coth":
				return	mul(
									neg(fdash),
									powen(
										apply("cosech", f),
										"2")), nil

			case "asinh", "arsinh", "arcsinh":
				return	div(
								fdash,
								apply("sqrt",
									addn(
										"1",
										powen(
											f,
											"2")))), nil

			case "acosh", "arcosh", "arccosh":
				return	div(
									fdash,
									mul(
										apply("sqrt",
											suben(
												f, "1")),
										apply("sqrt",
											addn(
												"1",
												f)))), nil

			case "atanh", "artanh", "arctanh":
				return	div(
									fdash,
									subne(
										"1",
										powen(
											f,
											"2"))), nil

			case "asech", "arsech", "arcsech":
				return	div(
									mul(
										fdash,
										apply("sqrt",
											div(
												subne(
													"1",
													f),
												addn(
													"1",
													f)))),
									mul(
										f,
										suben(
											f,
											"1"))), nil

			case "acsch", "arcsch", "arccsch", "acosech", "arcosech", "arccosech":
				return	div(
									neg(fdash),
									mul(
										apply("sqrt",
											addn(
												"1",
												divne(
													"1",
													powen(
														f,
														"2")))),
										powen(
											f,
											"2"))), nil

			case "acoth", "arcoth", "arccoth":
				return	div(
									fdash,
									subne(
										"1",
										powen(
											f,
											"2"))), nil
		}

	default:
		return nil, errors.New("Unknown operator: " + exp.Op)
	}

	return exp, nil
}

func add(e1, e2 *Expression) *Expression {
	return &Expression{"+", OP_LOW, e1, e2}
}

func addn(n string, e2 *Expression) *Expression {
	return add(&Expression{n, NUMBER, nil, nil}, e2)
}

func sub(e1, e2 *Expression) *Expression {
	return &Expression{"-", OP_LOW, e1, e2}
}

func subne(n string, e2 *Expression) *Expression {
	return sub(&Expression{n, NUMBER, nil, nil}, e2)
}

func suben(e1 *Expression, n string) *Expression {
	return sub(e1, &Expression{n, NUMBER, nil, nil})
}

func mul(e1, e2 *Expression) *Expression {
	return &Expression{"*", OP_MED, e1, e2}
}

func muln(n string, e2 *Expression) *Expression {
	return mul(&Expression{n, NUMBER, nil, nil}, e2)
}

func div(e1, e2 *Expression) *Expression {
	return &Expression{"/", OP_MED, e1, e2}
}

func divne(n string, e2 *Expression) *Expression {
	return div(&Expression{n, NUMBER, nil, nil}, e2)
}

func diven(e1 *Expression, n string) *Expression {
	return div(e1, &Expression{n, NUMBER, nil, nil})
}

func pow(e1, e2 *Expression) *Expression {
	return &Expression{"^", OP_HIGH, e1, e2}
}

func powne(n string, e2 *Expression) *Expression {
	return pow(&Expression{n, NUMBER, nil, nil}, e2)
}

func powen(e1 *Expression, n string) *Expression {
	return pow(e1, &Expression{n, NUMBER, nil, nil})
}

func apply(fn string, e *Expression) *Expression {
	return &Expression{fn, FUNC_PREFIX, e, nil}
}

func neg(e *Expression) *Expression {
	return muln("-1", e)
}

func no(num string) *Expression {
	return &Expression{num, NUMBER, nil, nil}
}