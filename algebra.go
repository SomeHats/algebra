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

			//fmt.Println("d/dx of expression: ", exp.Differentiate("x"))
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
