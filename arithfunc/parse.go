//Package arithfunc provides a function for converting a string based arithmetic function to a func variable.
package arithfunc

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
)

//Node types
const (
	operatorNode = iota
	variableNode = iota
	constantNode = iota
	functionNode = iota
)

//Operator types
const (
	addOp = iota
	subOp = iota
	mulOp = iota
	divOp = iota
	powOp = iota
)

//Operator symbols, the order must match the operator type order
var operators = []byte{
	addOp: '+',
	subOp: '-',
	mulOp: '*',
	divOp: '/',
	powOp: '^',
}

//Function types
const (
	absFunc  = iota
	sinFunc  = iota
	cosFunc  = iota
	tanFunc  = iota
	lnFunc   = iota
	logFunc  = iota
	asinFunc = iota
	acosFunc = iota
	atanFunc = iota
)

//Function definitions, the order must match the function type order
var functions = []string{
	absFunc:  "abs(",
	sinFunc:  "sin(",
	cosFunc:  "cos(",
	tanFunc:  "tan(",
	lnFunc:   "ln(",
	logFunc:  "log(",
	asinFunc: "asin(",
	acosFunc: "acos(",
	atanFunc: "atan(",
}

//Pre-defined constants
var constantMap = map[string]float64{
	"e":   math.E,
	"pi":  math.Pi,
	"phi": math.Phi,
}

//Struct containing a node which can contain either an operator, a variable, or a constant
type node struct {
	nodeType      int //Defines node type
	children      []*node
	operatorType  int
	functionType  int
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
	if root == nil {
		return nil, errors.New("Input function string is empty.")
	}

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
	//If node is nil, return zero. This will properly execute a negation
	if n == nil {
		return 0
	}

	//Check node type
	switch n.nodeType {
	case operatorNode:
		//If operator, recursively calculate children and execute operation
		switch n.operatorType {
		case addOp:
			return traverseAndCalc(n.children[0], vl...) + traverseAndCalc(n.children[1], vl...)
		case subOp:
			return traverseAndCalc(n.children[0], vl...) - traverseAndCalc(n.children[1], vl...)
		case mulOp:
			return traverseAndCalc(n.children[0], vl...) * traverseAndCalc(n.children[1], vl...)
		case divOp:
			return traverseAndCalc(n.children[0], vl...) / traverseAndCalc(n.children[1], vl...)
		case powOp:
			return math.Pow(traverseAndCalc(n.children[0], vl...), traverseAndCalc(n.children[1], vl...))
		}
	case functionNode:
		//If function, execute function on inner right tree
		switch n.functionType {
		case absFunc:
			return math.Abs(traverseAndCalc(n.children[0], vl...))
		case sinFunc:
			return math.Sin(traverseAndCalc(n.children[0], vl...))
		case cosFunc:
			return math.Cos(traverseAndCalc(n.children[0], vl...))
		case tanFunc:
			return math.Tan(traverseAndCalc(n.children[0], vl...))
		case lnFunc:
			return math.Log(traverseAndCalc(n.children[0], vl...))
		case logFunc:
			return math.Log10(traverseAndCalc(n.children[0], vl...))
		case asinFunc:
			return math.Asin(traverseAndCalc(n.children[0], vl...))
		case acosFunc:
			return math.Acos(traverseAndCalc(n.children[0], vl...))
		case atanFunc:
			return math.Atan(traverseAndCalc(n.children[0], vl...))
		}
	case constantNode:
		return n.constantValue
	case variableNode:
		if n.variableIndex >= len(vl) {
			panic(errors.New(fmt.Sprintf("Variable with index %d is required. Not enough input variables were provided in function call.", n.variableIndex)))
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
	for k, operator := range operators {
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
				//If the operator is - it might be a negation operation, in this case there should be another operator to the left of it, check for it
			charLoop:
				for i2 := i; i2 >= 0; i2-- {
					//Check the rune at the current position
					switch line[i2] {
					case '-':
						i = i2
					case '+':
						fallthrough
					case '*':
						fallthrough
					case '/':
						fallthrough
					case '^':
						i = i2 //Reposition i to this operator and break out of for loop
						break charLoop
					case ' ':
						//Do nothing - keep checking
					default:
						break charLoop //if any value that isn't an operator or space is detected, the last operator has been found
					}
				}

				//Create and return an operator node
				left := createNode(line[:i])
				right := createNode(line[i+1:])

				//Right child being nil is invalid in all cases, left child being nil is only valid for - operator
				if right == nil || (value != '-' && left == nil) {
					panic(errors.New(fmt.Sprintf("An operator with symbol %s is missing a child on which to operate. Check that the function is properly formatted.", string(value))))
				}

				return &node{
					nodeType:     operatorNode,
					operatorType: k,
					children:     []*node{left, right},
				}
			}
		}
	}

	//If no operators have been found, the remaining value represents either a variable, a constant, a function, or following a negation operator
	//Return nil node if line is empty, this is only valid in the case of a negation operation, will cause a fault for everything else
	if len(line) == 0 {
		return nil
	}

	//Attempt to parse for a numerical constant second. Check for e prevents values defined in scientific notation such as 1e5 to go through without error.
	//This behavior is enforced because 1e-5 would throw an error anyway due to the negative sign.
	num, err := strconv.ParseFloat(line, 64)
	if err == nil && !strings.ContainsAny(line, "eE") {
		//Parsing constant succeeded, create and return constant node
		return &node{
			nodeType:      constantNode,
			constantValue: num,
		}
	}

	//Check if value is a valid pre-defined constant third
	num, ok := constantMap[line]
	if ok {
		//Fetching constant succeeded, create and return constant node
		return &node{
			nodeType:      constantNode,
			constantValue: num,
		}
	}

	//Check if line contains a valid function fourth
	for k, v := range functions {
		if strings.HasPrefix(line, v) {
			return &node{
				nodeType:     functionNode,
				functionType: k,
				children:     []*node{createNode(line[len(v) : len(line)-1])},
			}
		}
	}

	//Attempt to parse for a variable last
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
		nodeType:      variableNode,
		variableIndex: int(varIdx),
	}
}
