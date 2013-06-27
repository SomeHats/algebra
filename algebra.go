package main

import (
	"./algebra"
	"bufio"
	"fmt"
	"os"
)

func main() {

	for {
		fmt.Printf("Enter an expression: ")

		expStr := readLine()

		fmt.Println("Expression tree:")
		exp, err := algebra.Parse(expStr)
		if err == nil {
			fmt.Println(exp)
			//fmt.Println(exp.ToLatex())

			fmt.Println("Simplified:")
			//exp = exp.Simplify()
			//fmt.Println(exp)
			//fmt.Println(exp.UnTree())

			exp, err := exp.Differentiate("x")

			if err == nil {
				fmt.Println("d/dx of expression: \n", exp)
				fmt.Println("Simplified: \n", exp.Simplify().UnTree())
			} else {
				fmt.Println("Error: ", err)
			}

		} else {
			fmt.Println("Error: ", err)
		}
		//return
	}
}

func readLine() string {
	bytes, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	return string(bytes)
}
