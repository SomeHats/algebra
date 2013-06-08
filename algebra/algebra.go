package algebra

import (
	"strings"
)

const (
	NUMBER       uint8 = 0
	FUNC_PREFIX        = 1
	CONSTANT           = 2
	VARIABLE           = 3
	PAREN_OPEN         = 4
	PAREN_CLOSE        = 5
	FUNC_POSTFIX       = 6
	OP_LOW             = 7
	OP_MED             = 8
	OP_HIGH            = 9
	EQUALS             = 10
)

type Expression struct {
	Op    string
	Type  uint8
	Left  *Expression
	Right *Expression
}

func (e *Expression) String() string {
	out, str := e.Op, ""

	if e.Left != nil {
		str = e.Left.String()
		if e.Right == nil {
			str = strings.Replace(str, "\n", "\n   ", -1)
			out += "\n└─ " + str
		} else {
			str = strings.Replace(str, "\n", "\n│  ", -1)
			out += "\n├─ " + str
		}
	}

	if e.Right != nil {
		str = e.Right.String()
		str = strings.Replace(str, "\n", "\n   ", -1)
		out += "\n└─ " + str
	}

	return out
}

func (e *Expression) UnTree() string {
	switch e.Type {
	case NUMBER, VARIABLE, CONSTANT, EQUALS:
		return e.Op

	case OP_LOW, OP_MED, OP_HIGH:
		return "(" + e.Left.UnTree() + " " + e.Op + " " + e.Right.UnTree() + ")"

	case FUNC_PREFIX:
		return e.Op + "(" + e.Left.UnTree() + ")"
	}

	return "CAN'T UNTREE: " + e.Op
}
