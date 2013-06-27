package algebra

import (
	"errors"
	"regexp"
	"strings"
	//"fmt"
)

// Token types:
// 0 - number
// 1 - prefix function (sin, cos, tan, ln, etc)
// 2 - constants (e, pi)
// 3 - variable
// 4 - open parenthesis
// 5 - close parenthesis
// 6 - postfix function (!)
// 7 - Low precedence operator (+, -)
// 8 - Medium precedence operator (*, /)
// 9 - High precedence operator (^)
// 10 - Equals

type token struct {
	Type  uint8
	Value string
}

var precedence = []uint8{EQUALS, OP_LOW, OP_MED, OP_HIGH}

func Parse(s string) (*Expression, error) {
	s = strings.ToLower(s)

	tokens := tokenize(s)

	// Check that brackets are balanced:
	depth := 0
	for _, token := range tokens {
		if token.Type == PAREN_OPEN {
			depth++
		} else if token.Type == PAREN_CLOSE {
			depth--
		}
	}

	if depth != 0 {
		if depth < 0 {
			return &Expression{}, errors.New("Unmatched ')'")
		} else {
			return &Expression{}, errors.New("Unmatched '('")
		}
	}

	// Parse the thing!
	return parseTokens(tokens)
}

func parseTokens(tokens []token) (*Expression, error) {
	// parseTokens will parse one operator, then call itself on the
	// operands

	tokens = debracketise(tokens)

	// Just 1 token is probably a value on its own
	if len(tokens) == 1 {
		return &Expression{tokens[0].Value, tokens[0].Type, nil, nil}, nil
	}

	// search for operators outside of brackets
	for _, op := range precedence {
		depth := 0
		for i := len(tokens) - 1; i >= 0; i-- {
			skip := false

			if tokens[i].Type == PAREN_CLOSE {
				depth++
			} else if tokens[i].Type == PAREN_OPEN {
				depth--
			}

			// Handle implicit multiplication:
			if op == 8 && i != 0 && depth == 0 {
				prevType := tokens[i-1].Type
				currentType := tokens[i].Type

				switch prevType {
				case NUMBER, CONSTANT, VARIABLE, PAREN_CLOSE: // implicit multiplication
					switch currentType {
					case NUMBER, FUNC_PREFIX, CONSTANT, VARIABLE, PAREN_OPEN:
						left, err := parseTokens(tokens[:i])
						if err != nil {
							return nil, err
						}
						right, err := parseTokens(tokens[i:])
						if err != nil {
							return nil, err
						}
						return &Expression{"*", OP_MED, left, right}, nil
					}
				}
			}

			// Process operators
			if (tokens[i].Type == op) && (depth == 0) {
				var (
					left *Expression
					err  error
				)
				// Handle unary '-'
				if (op == 7) && (tokens[i].Value == "-") && ((i == 0) || (tokens[i-1].Type >= 7)) {
					left = &Expression{"0", NUMBER, nil, nil}
					if i != 0 {
						skip = true
					}
				} else {
					left, err = parseTokens(tokens[:i])
					if err != nil {
						return nil, err
					}
				}

				if !skip {
					right, err := parseTokens(tokens[i+1:])
					if err != nil {
						return nil, err
					}

					return &Expression{tokens[i].Value, tokens[i].Type, left, right}, nil
				}
			}
		}
	}

	// prefix & postfix functions
	for i := 1; i < len(tokens); i++ {
		currentType := tokens[i].Type
		prevType := tokens[i-1].Type

		if prevType == FUNC_PREFIX {
			left, err := parseTokens(tokens[i:])
			if err != nil {
				return nil, err
			}
			return &Expression{tokens[i-1].Value, prevType, left, nil}, nil
		}

		// postfix functions:
		if currentType == FUNC_POSTFIX {
			left, err := parseTokens(tokens[:i])
			if err != nil {
				return nil, err
			}
			return &Expression{tokens[i].Value, currentType, left, nil}, nil
		}
	}

	return nil, errors.New("Couldn't parse tokens")
}

func debracketise(tokens []token) []token {
	if len(tokens) < 3 {
		return tokens
	}
	depth := 0
	hitsZero := false
	for i, op := range tokens {
		if op.Type == PAREN_OPEN {
			depth++
		} else if op.Type == PAREN_CLOSE {
			depth--
		}

		if (depth == 0) && (i != len(tokens)-1) {
			hitsZero = true
		}
	}

	if hitsZero {
		return tokens
	}

	return tokens[1 : len(tokens)-1]
}

func tokenize(s string) []token {
	matchers := []*regexp.Regexp{
		regexp.MustCompile("([0-9]+)(\\.[0-9]+)?(e-?[0-9]+)?"),
		regexp.MustCompile("((a(rc?)?)?(cosec|sin|cos|tan|sec|csc|cot)h?)|ln|log|sqrt"),
		regexp.MustCompile("e|i|pi"),
		regexp.MustCompile("[a-z]"),
		regexp.MustCompile("\\("),
		regexp.MustCompile("\\)"),
		regexp.MustCompile("!"),
		regexp.MustCompile("[+-]"),
		regexp.MustCompile("[*/]"),
		regexp.MustCompile("\\^"),
		regexp.MustCompile("=")}

	// order shows the order that tokens come in. -1 means that spot
	// is empty, any other shows its position in the tokens array
	order := make([]int, len(s))

	for i := 0; i < cap(order); i++ {
		order[i] = -1
	}

	// List of tokens
	tokens := make([]token, 0, 10)

	// List of characters from the string that have been tokenized
	tokenized := make([]bool, len(s))

	//----------THE ACTUAL TOKENIZING BIT. WOOOOO!-----------//

	for i, match := range matchers {
		// Find positions of all the things
		nos := match.FindAllStringIndex(s, -1)

		// Mark them as found and store their values
		for _, n := range nos {
			ignore := false
			for i := n[0]; i < n[1]; i++ {
				if tokenized[i] {
					ignore = true
				}
				tokenized[i] = true
			}
			if !ignore {
				order[n[0]] = len(tokens)
				tokens = append(tokens, token{uint8(i), s[n[0]:n[1]]})
			}
		}
	}

	out := make([]token, 0, len(tokens))
	for i := 0; i < len(s); i++ {
		if order[i] != -1 {
			out = append(out, tokens[order[i]])
		}
	}

	return out
}

func (e *Expression) ToLatex() string {
	opType, op := e.Type, e.Op

	switch opType {
	case NUMBER, CONSTANT, VARIABLE:
		if op == "pi" {
			return " \\pi "
		}

		return op

	case FUNC_PREFIX:
		var arg string
		switch e.Left.Type {
			case NUMBER, CONSTANT, VARIABLE, FUNC_PREFIX:
				arg = " " + e.Left.ToLatex() + " "

			default:
				arg = " \\left ( " + e.Left.ToLatex() + " \\right ) "
		}
		switch op{
			case "ln", "sin", "cos", "tan", "sec", "csc", "cot",
			"sinh", "cosh", "tanh", "sech", "csch", "coth":
				return " \\" + op + arg

			case "cosec":
				return " \\csc" + arg

			case "asin", "arsin", "arcsin":
				return " \\arcsin" + arg

			case "acos", "arcos", "arccos":
				return " \\arccos" + arg

			case "atan", "artan", "arctan":
				return " \\arctan" + arg

			case "asec", "arsec", "arcsec":
				return " \\arcsec" + arg

			case "acsc", "arcsc", "arccsc", "acosec", "arcosec", "arccosec":
				return " \\arccsc" + arg

			case "acot", "arcot", "arccot":
				return " \\arccot" + arg

			case "log":
				return " \\log_{10} " + arg

			case "sqrt":
				return " \\sqrt{" + e.Left.ToLatex() + "} "
		}

	case FUNC_POSTFIX:
		var arg string
		switch e.Left.Type {
			case NUMBER, CONSTANT, VARIABLE, FUNC_PREFIX:
				arg = " " + e.Left.ToLatex() + " "

			default:
				arg = " \\left ( " + e.Left.ToLatex() + " \\right ) "
		}

		return arg + op + " "

	case OP_LOW, OP_MED, OP_HIGH:
		var left, right string

		if op == "/" {
			left = " " + e.Left.ToLatex() + " "
			right = " " + e.Right.ToLatex() + " "
			return " \\frac{" + left + "}{" + right + "}"
		}

		switch e.Left.Type {
			case NUMBER, CONSTANT, VARIABLE, FUNC_PREFIX:
				left = " " + e.Left.ToLatex() + " "

			case OP_LOW, OP_MED, OP_HIGH:
				if e.Left.Type < opType {
					left = " \\left ( " + e.Left.ToLatex() + " \\right ) "
				} else {
					left = " " + e.Left.ToLatex() + " "
				}

			default:
				left = " \\left ( " + e.Left.ToLatex() + " \\right ) "
		}
		switch e.Right.Type {
			case NUMBER, CONSTANT, VARIABLE, FUNC_PREFIX:
				right = " " + e.Right.ToLatex() + " "

			case OP_LOW, OP_MED, OP_HIGH:
				if e.Right.Type < opType {
					right = " \\left ( " + e.Right.ToLatex() + " \\right ) "
				} else {
					right = " " + e.Right.ToLatex() + " "
				}

			default:
				right = " \\left ( " + e.Right.ToLatex() + " \\right ) "
		}

		if op == "*" {
			op = ""
		}

		return left + op + right
	}

	return "MISSING: " + op
}
