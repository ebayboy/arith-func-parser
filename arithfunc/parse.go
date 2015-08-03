package arithfunc

import (
	"errors"
	"math"
	"runtime"
	"strconv"
	"strings"
)

const (
	operatorNode = iota
	variableNode = iota
	constantNode = iota
)

//Struct containing a node which can contain either an operator, a variable, or a constant
type node struct {
	left  *node
	right *node
	value string

	nodeType      int
	constantValue float64
	variableIndex int
}

//Parse takes in a string denoting an arithmetic function with optional variable values formatted as V0, V1, etc.
//Returns a func that will evalute the function for the provided variable values.
func Parse(s string) (result func(vl ...float64) (float64, error), err error) {
	//Because the called functions are recursive, in order to return an error the function will panic internally
	//the panic will be recovered here and returned as an error
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	root := createNode(s)

	return func(vl ...float64) (value float64, err error) {
		//Set up recovery handling on recursive function call
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(runtime.Error); ok {
					panic(r)
				}
				err = r.(error)
			}
		}()

		//Create var map to map variable values to corresponding string values that were originally passed to parse function
		return traverseAndCalc(root, vl...), nil
	}, nil
}

func traverseAndCalc(n *node, vl ...float64) float64 {
	switch n.nodeType {
	case operatorNode:
		switch n.value {
		case "+":
			return traverseAndCalc(n.left, vl...) + traverseAndCalc(n.right, vl...)
		case "-":
			return traverseAndCalc(n.left, vl...) - traverseAndCalc(n.right, vl...)
		case "*":
			return traverseAndCalc(n.left, vl...) * traverseAndCalc(n.right, vl...)
		case "/":
			return traverseAndCalc(n.left, vl...) / traverseAndCalc(n.right, vl...)
		case "^":
			return math.Pow(traverseAndCalc(n.left, vl...), traverseAndCalc(n.right, vl...))
		}
	case constantNode:
		return n.constantValue
	case variableNode:
		if n.variableIndex >= len(vl) {
			panic(errors.New("Variable with index %d is required. Not enough input variables were provided in function call."))
		}
		return vl[n.variableIndex]
	}

	panic("Invalid node type.")
}

func isSurroundedByMatchingParentheses(line string) bool {
	//Only check for closing parenthesis if line is actually surrounded by parentheses
	if len(line) < 2 || !strings.HasPrefix(line, "(") || !strings.HasSuffix(line, ")") {
		return false
	}

	//Search for closing parenthesis that is not the last parenthesis in the inner string
	search := line[1 : len(line)-1]

	//Count parentheses, the count should never go below 1 else we have found a closing parenthesis
	count := 1
	for _, runeValue := range search {
		//Increment of decrement count as parentheses are found
		switch runeValue {
		case '(':
			count++
		case ')':
			count--
		}

		//Matching closing parenthesis found before the last parenthesis
		if count == 0 {
			return false
		}
	}

	return true
}

var operators = []byte{'+', '-', '*', '/', '^'}

//Creates a node from the given string
func createNode(line string) *node {
	//Trim white space if there is any
	line = strings.TrimSpace(line)

	//While line is surrounded by matching parentheses, remove them
	for isSurroundedByMatchingParentheses(line) {
		line = strings.TrimSpace(line[1 : len(line)-1])
	}

	//Iterate through characters looking for operation characters with respect to order of operations.
	//Anything in parentheses is ignored as it will be handled later
	for _, operator := range operators {
		parenthCount := 0

		//Iterate through string backwards to preserve left to right operation
		for i := len(line) - 1; i >= 0; i-- {
			value := line[i]

			switch {
			case value == ')':
				parenthCount++
			case value == '(':
				parenthCount--
			case parenthCount == 0 && value == operator:
				//Create and return an operator node
				return &node{
					value:    string(value),
					nodeType: operatorNode,
					left:     createNode(line[:i]),
					right:    createNode(line[i+1:]),
				}
			}
		}
	}

	//If no operators have been found, the remaining value represents either a variable or a constant
	//Attempt to parse for a constant first
	num, err := strconv.ParseFloat(line, 64)
	if err != nil {
		recoverableError := errors.New("The function defined by the string is improperly formatted. Only variables of the form V# (V0 for variable 0) are allowed. " +
			"Negative constants cannot be defined as this will get confused with the subtract operator, define these as (0 - 5) for -5. " +
			"Ensure all parentheses are matching pairs.")

		//Parsing for a constant failed, either the value is a variable of the function was improperly formatted
		//Attempt to parse for a variable
		if len(line) < 2 || strings.ToUpper(line[:1]) != "V" {
			panic(recoverableError)
		}

		//Parse for variable index
		varIdx, err := strconv.ParseUint(line[1:], 10, 32)
		if err != nil {
			panic(recoverableError)
		}

		//Variable parse worked
		return &node{
			value:         line,
			nodeType:      variableNode,
			variableIndex: int(varIdx),
		}
	}

	//Parsing constant succeeded, create and return constant node
	return &node{
		value:         line,
		nodeType:      constantNode,
		constantValue: num,
	}
}
