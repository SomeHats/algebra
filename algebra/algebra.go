package algebra

import (
	"strings"
)

type Expression struct {
	Op    string
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
