package arithfunc

import (
	"math"
	"testing"
)

func init() {
	RegisterFunction("sum3", 3, func(vars ...float64) float64 {
		return vars[0] + vars[1] + vars[2]
	})
}

func TestParse(t *testing.T) {
	inputs := []float64{1, -0.5}

	testParseHelper(t, 2, "(V0 * 4)^(1/2)", inputs...)
	testParseHelper(t, 1, " ( V0 - -(sin(pi/2))   + -abs(cos(pi)))", 1)
	testParseHelper(t, 5, "5")
	testParseHelper(t, 36, "36")
	testParseHelper(t, 1.5, "1.5")
	testParseHelper(t, -0.5, "V0", -0.5)
	testParseHelper(t, 0.5, "V0 / V0 + V1", inputs...)
	testParseHelper(t, 2, "V0 / (V0 + V1)", inputs...)
	testParseHelper(t, -1, "V0 - V0 - V0", inputs...)
	testParseHelper(t, 4, "V0 / V1 / V1", inputs...)
	testParseHelper(t, -2, "(0 - 1) / (V0 + V1)", inputs...)
	testParseHelper(t, -2, "( (0 - 1) ) / (V0 + V1)", inputs...)
	testParseHelper(t, 4, "(((V0001 - V00)/1.5) + 5)", inputs...)
	testParseHelper(t, 16, "(V0 * 4)^2", inputs...)
	testParseHelper(t, 9, "6 * 1 / 2 * 3")
	testParseHelper(t, -12, "1 + 2 / 4 * 5 + 1 / 2 - (7 / 2 * 2 + 9)")
	testParseHelper(t, 10, "5 - (-5 - -5) + V0", 5)
	testParseHelper(t, 10, "5--5")
	testParseHelper(t, 0, "-5--5")
	testParseHelper(t, math.Inf(0), "1 / 0")
	testParseHelper(t, math.Inf(-1), "(0 - 1) / 0")
	testParseHelper(t, math.Pi, "pi")
	testParseHelper(t, 1, "abs(-1)")
	testParseHelper(t, 1, "abs(1)")
	testParseHelper(t, 1, "cos(0)")
	testParseHelper(t, 0, "cos(pi/2)")
	testParseHelper(t, 0, "sin(0)")
	testParseHelper(t, 1, "sin(pi/2)")
	testParseHelper(t, 5, "log(10^5)")
	testParseHelper(t, 5, "ln(e^5)")
	testParseHelper(t, 1e-5, "10^(-5)")
	testParseHelper(t, math.Pi/4, "asin(sin(pi/4))")
	testParseHelper(t, math.Pi/4, "acos(cos(pi/4))")
	testParseHelper(t, math.Pi/4, "atan(tan(pi/4))")
	testParseHelper(t, -math.Pi/4, "atan2(-1, 1)")
	testParseHelper(t, 2, "sqrt(4)")
	testParseHelper(t, 10, "sum3(10/5, 3, 2.5 * 2)") //test parsing a custom function
}

func testParseHelper(t *testing.T, answer float64, fStr string, vl ...float64) {
	f, err := Parse(fStr)
	if err != nil {
		t.Logf("Failed on input: %v", fStr)
		t.Logf(err.Error())
		t.Fail()
		return
	}

	result, err := f(vl...)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
		return
	}

	if math.Abs(result-answer) > 1e-10 {
		t.Logf("Answer is incorrect. %s = %f (%v) Answer: %f\r\n", fStr, result, vl, answer)
		t.Fail()
		return
	}

	//fmt.Printf("Answer is correct. %s = %f (%v) Answer: %f\r\n", fStr, result, vl, answer)
}

func TestParseErrors(t *testing.T) {
	testParseErrorsHelper(t, "V0 / (V0 + V1")
	testParseErrorsHelper(t, "A0 / 1")
	testParseErrorsHelper(t, "()())((()))))))")
	testParseErrorsHelper(t, "5--5-")
	testParseErrorsHelper(t, "5 - *V0")
	testParseErrorsHelper(t, "5 - *3")
	testParseErrorsHelper(t, "5*")
	testParseErrorsHelper(t, "/4")
	testParseErrorsHelper(t, "")
	testParseErrorsHelper(t, " ")
	testParseErrorsHelper(t, "abc")
	testParseErrorsHelper(t, "⌘")
	testParseErrorsHelper(t, "1e5")
	testParseErrorsHelper(t, "1e-5")
	testParseErrorsHelper(t, "1E5")
	testParseErrorsHelper(t, "1E-5")
	testParseErrorsHelper(t, "sin()")
	testParseErrorsHelper(t, "abs()")
	testParseErrorsHelper(t, "abs(5")
	testParseErrorsHelper(t, "abs(")
	testParseErrorsHelper(t, "abs((5/2.5)")
	testParseErrorsHelper(t, "abs(5, 7)")
	testParseErrorsHelper(t, "atan2(1)")
	testParseErrorsHelper(t, "fake(1, 2, 3)")
	testParseErrorsHelper(t, "sum3(1,2,3,4)") //sum3 was registered before, test to make sure it errors out on providing the wrong number of variables
}

func testParseErrorsHelper(t *testing.T, fStr string) {
	_, err := Parse(fStr)
	if err == nil {
		t.Logf("Failed to return error parsing: %s", fStr)
		t.Fail()
		return
	}

	//	//If error is returned, that is a pass
	//	fmt.Printf("Success! %s produced error: %s\r\n", fStr, err.Error())
}

func TestExecutionErrors(t *testing.T) {
	testExecutionErrorsHelper(t, "V0 + V1 + V2", 1, 2)
	testExecutionErrorsHelper(t, "V9 + 5", 1, 2, 3, 4)
}

func testExecutionErrorsHelper(t *testing.T, fStr string, vl ...float64) {
	f, err := Parse(fStr)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
		return
	}

	_, err = f(vl...)
	if err == nil {
		t.Logf("Failed to return error executing: %s (%v)", fStr, vl)
		t.Fail()
		return
	}
}
