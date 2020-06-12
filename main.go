package main

import (
	"fmt"
	"github.com/ebayboy/arith-func-parser/arithfunc"
)

func main() {
	f, err := arithfunc.Parse("abs(-5 - V0) / V0 + V1^(1/2)")
	if err != nil {
		panic("Encountered error trying to parse string.")
	}

	result, err := f(1, 4)
	if err != nil {
		panic("Encountered error trying to execute function.")
	}

	fmt.Println(result)
	//Result should be 8
}
