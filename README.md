# arith-func-parser
Go package for parsing an arithmetic function with optional variables. Creates a func variable that can be quickly evaluated given a matching set of variable values.

## Table of content

- [Features](#features)
- [Install](#install)
- [Examples](#examples)
- [Limitations](#limitations)

## Features

- Pass arithmetic function as string and return func variable that can be executed to obtain the result.
- Supports any number of variables in the function definition which can then be evaluated using values passed to the returned function.

## Install

This package is "go-gettable":

    go get github.com/JLaferri/arith-func-parser/arithfunc

## Examples

``` go
package main

import (
	"fmt"
	"github.com/JLaferri/arith-func-parser/arithfunc"
)

func main() {
	f, err := arithfunc.Parse("-5 - V0 / V0 + V1^(1/2)")
	if err != nil {
		panic("Encountered error trying to parse string.")
	}

	result, err := f(1, 4)
	if err != nil {
		panic("Encountered error trying to execute fucntion.")
	}

	fmt.Println(result)
}
```

## Limitations

- Does not support complex numbers.
- The only valid operators are (+, -, *, /, and ^)
