package algebra

import (
  "math/big"
  "strings"
)

func (exp *Expression) Simplify() *Expression {
  if exp.Left != nil {
    exp.Left = exp.Left.Simplify()
  }
  if exp.Right != nil {
    exp.Right = exp.Right.Simplify()
  }

  if exp.IsConstant() {
    exp = exp.GetConstantTree()
  }

  left := exp.Left
  right := exp.Right

  op := exp.Op

  switch op {
    case "+":
      if left.Type == NUMBER && left.Op == "0" {
        exp = right
      } else if right.Type == NUMBER && right.Op == "0" {
        exp = left
      }

    case "-":
      if right.Type == NUMBER && right.Op == "0" {
        exp = left
      }

    case "*":
      if left.Type == NUMBER && (left.Op == "0" || left.Op == "1") {
        if left.Op == "0" {
          exp = left
        } else {
          exp = right
        }
      } else if right.Type == NUMBER && (right.Op == "0" || right.Op == "1") {
        if right.Op == "0" {
          exp = right
        } else {
          exp = left
        }
      }

    case "/":
      if left.Type == NUMBER && left.Op == "0" {
        exp = left
      }

    case "^":
      if right.Type == NUMBER && right.Op == "1" {
        exp = left
      }
  }

  return exp
}

func (exp *Expression) IsConstant() bool {
  if exp.Type == VARIABLE {
    return false
  }

  if exp.Left == nil && exp.Right == nil {
    return true
  }

  if exp.Left == nil {
    return exp.Right.IsConstant()
  }

  if exp.Right == nil {
    return exp.Left.IsConstant()
  }

  return exp.Left.IsConstant() && exp.Right.IsConstant()
}

func (exp *Expression) GetConstantTree() *Expression {
  if exp.Left == nil || exp.Right == nil {
    if exp.Type == NUMBER {
      n, err := big.NewRat(1, 1).SetString(exp.Op)
      if !err {
        panic("Can't parse number: " + exp.Op)
      }

      return ratToExp(n)
    }
    return exp
  }

  exp.Left = exp.Left.GetConstantTree()
  exp.Right = exp.Right.GetConstantTree()

  left := exp.Left
  right := exp.Right

  op := exp.Op
  opType := exp.Type

  switch opType {
    case OP_LOW, OP_MED, OP_HIGH:
      if left.Type == NUMBER && right.Type == NUMBER {
        n1, err := big.NewRat(1, 1).SetString(left.Op)
        if !err {
          panic("Can't parse number: " + left.Op)
        }
        n2, err := big.NewRat(1, 1).SetString(right.Op)
        if !err {
          panic("Can't parse number: " + left.Op)
        }

        switch op {
          case "+":
            n1.Add(n1, n2)

          case "-":
            n1.Sub(n1, n2)

          case "*":
            n1.Mul(n1, n2)

          case "/":
            n1.Quo(n1, n2)

          case "^":
            if n2.IsInt() {
              n := n2.Num()
              a := n1.Num()
              b := n1.Denom()

              n1.SetFrac(a.Exp(a, n, big.NewInt(0)), b.Exp(b, n, big.NewInt(0)))
            } else {
              return exp
            }
        }

        return ratToExp(n1)
      }
  }

  return exp
}

func ratToExp(n *big.Rat) *Expression {
  if n.IsInt() {
    return &Expression{n.Num().String(), NUMBER, nil, nil}
  }
  return  &Expression{"/", OP_MED,
            &Expression{n.Num().String(), NUMBER, nil, nil},
            &Expression{n.Denom().String(), NUMBER, nil, nil}}
}

func isInt (s string) bool {
  return !strings.Contains(s, ".")
}
