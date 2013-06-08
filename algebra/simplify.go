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
    return exp
  }

  exp.Left = exp.Left.GetConstantTree()
  exp.Right = exp.Right.GetConstantTree()

  left := exp.Left
  right := exp.Right

  op := exp.Op

  switch op {
    case "+", "-":
      if left.Type == NUMBER && right.Type == NUMBER {
        if isInt(left.Op) && isInt(right.Op) {
          n1, err := big.NewInt(0).SetString(left.Op, 10)
          if !err {
            panic("Can't parse int: " + left.Op)
          }
          n2, err := big.NewInt(0).SetString(right.Op, 10)
          if !err {
            panic("Can't parse int: " + right.Op)
          }

          if op == "+" {
            n1.Add(n1, n2)
          } else {
            n1.Sub(n1, n2)
          }

          return &Expression{n1.String(), NUMBER, nil, nil}
        }
      }
  }

  return exp
}

func isInt (s string) bool {
  return !strings.Contains(s, ".")
}
