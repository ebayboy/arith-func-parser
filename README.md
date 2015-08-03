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
	"github.com/JLaferri/arith-func-parser/arithfunc"
	"fmt"
)

func main() {
	f, err := Parse("5 - V0 / V0 + V1^(1/2)")
	if err != nil {
	  panic()
	}

	result, err := f(1, 4)
	if err != nil {
	  panic()
	}
	
	fmt.Println(result)
}
```

## Limitations

- Does not currently support back to back operator symbols. In other words an input of "8 - -8" is invalid. In this case "8 - (0 - 8)" would have to be provided to achieve the expected result.
- Does not support complex numbers.
- The only valid operators are (+, -, *, /, and ^)
