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
			return exp, nil
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
			exp.Left = fdash
			exp.Right = gdash

		case "*":
			exp.Left = &Expression{"*", OP_MED, fdash, g}
			exp.Right = &Expression{"*", OP_MED, f, gdash}
			exp.Op = "+"
			exp.Type = OP_LOW

		case "/":
			exp.Left = &Expression{"-", OP_LOW,
				&Expression{"*", OP_MED, fdash, g},
				&Expression{"*", OP_MED, f, gdash}}

			exp.Right = &Expression{"^", OP_HIGH, g,
				&Expression{"2", NUMBER, nil, nil}}

		case "^":
			exp.Op = "*"
			exp.Type = OP_MED

			// left: f(x)^(g(x)-1)
			exp.Left = &Expression{"^", OP_HIGH, f,
				&Expression{"-", OP_LOW, g, &Expression{"1", NUMBER, nil, nil}}}

			//right: (g(x)*f'(x))+(f(x)*ln(f(x))*g'(x))
			exp.Right = &Expression{"+", OP_LOW,
				&Expression{"*", OP_MED, g, fdash},
				&Expression{"*", OP_MED, f,
					&Expression{"*", OP_MED, gdash,
						&Expression{"ln", FUNC_PREFIX, f, nil}}}}
		}

	default:
		return nil, errors.New("Unknown operator: " + exp.Op)
	}

	return exp, nil
}
