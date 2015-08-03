package arithfunc

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

const (
	operator = iota
	variable = iota
	constant = iota
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
func Parse(s string) func(vl ...float64) (float64, error) {
	root := createNode(s)

	return func(vl ...float64) (float64, error) {
		//Create var map to map variable values to corresponding string values that were originally passed to parse function
		vm := map[string]float64{}

		if vl != nil {
			//Fill up var map with values
			for i := 0; i < len(vl); i++ {
				vm[fmt.Sprintf("V%d", i)]
			}
		}

		return traverseAndCalc(root, vm)
	}
}

func traverseAndCalc(n *node, vm map[string]float64) (float64, error) {
	switch n {
	case "+":
		return traverseAndCalc(n.left, vm) + traverseAndCalc(n.right, vm)
	case "-":
		return traverseAndCalc(n.left, vm) - traverseAndCalc(n.right, vm)
	case "*":
		return traverseAndCalc(n.left, vm) * traverseAndCalc(n.right, vm)
	case "/":
		return traverseAndCalc(n.left, vm) / traverseAndCalc(n.right, vm)
	case "^":
		return math.Pow(traverseAndCalc(n.left, vm), traverseAndCalc(n.right, vm))
	default:
		//When the node does not contain an operator, it is either a number, a variable, or invalid
		var num float64

		//First we will attempt to read it as a number
		if num, err = strconv.ParseFloat(n.value, 64); err != nil {
			//If reading the number failed, attempt to load the value for the assumed variable
			if num, ok := vm[n.value]; !ok {
				return 0, errors.New("Not enough variable values were provided for the quantity of variables in the parsed function.")
			}
		}

		return num, nil
	}
}

//Creates a node from the given string
func createNode(s string) *node {

}
