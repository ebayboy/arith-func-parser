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
- Supports function computations such as abs, sin, log, etc.
- Supports some built in constants such as pi and e.

## Supported Functions

Function  | Go Function
------- | -------
abs(x)  | math.Abs(x)
sin(x)  | math.Sin(x)
cos(x)  | math.Cos(x)
tan(x)  | math.Tan(x)
ln(x)	| math.Log(x)
log(x)	| math.Log10(x)
asin(x)	| math.Asin(x)
acos(x)	| math.Acos(x)
atan(x)	| math.Atan(x)
atan2(y, x) | math.Atan2(y, x)

## Supported Constants

Constant | Value
-------- | -------
e	| math.E
pi	| math.Pi
phi	| math.Phi

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
	f, err := arithfunc.Parse("abs(-5 - V0) / V0 + V1^(1/2)")
	if err != nil {
		panic("Encountered error trying to parse string.")
	}

	result, err := f(1, 4)
	if err != nil {
		panic("Encountered error trying to execute fucntion.")
	}

	fmt.Println(result)
	//Result should be 8
}
```

## Limitations

- Does not support complex numbers.
- The only valid operators are (+, -, *, /, and ^)
