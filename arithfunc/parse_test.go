package arithfunc

import (
	"math"
	"testing"
)

func TestParse(t *testing.T) {
	inputs := []float64{1, -0.5}

	testParseHelper(t, 5, "5")
	testParseHelper(t, 36, "36")
	testParseHelper(t, 1.5, "1.5")
	testParseHelper(t, -0.5, "V0", -0.5)
	testParseHelper(t, 0.5, "V0 / V0 + V1", inputs...)
	testParseHelper(t, 2, "V0 / (V0 + V1)", inputs...)
	testParseHelper(t, -1, "v0 - v0 - v0", inputs...)
	testParseHelper(t, 4, "V0 / V1 / V1", inputs...)
	testParseHelper(t, -2, "(0 - 1) / (V0 + V1)", inputs...)
	testParseHelper(t, -2, "( (0 - 1) ) / (V0 + V1)", inputs...)
	testParseHelper(t, 4, "(((V0001 - V00)/1.5) + 5)", inputs...)
	testParseHelper(t, 2, "(V0 * 4)^(1/2)", inputs...)
	testParseHelper(t, 16, "(V0 * 4)^2", inputs...)
	testParseHelper(t, 9, "6 * 1 / 2 * 3")
	testParseHelper(t, -12, "1 + 2 / 4 * 5 + 1 / 2 - (7 / 2 * 2 + 9)")
	testParseHelper(t, 10, "5 - (-5 - -5) + V0", 5)
	testParseHelper(t, 10, "5--5")
	testParseHelper(t, 0, "-5--5")
	testParseHelper(t, math.Inf(0), "1 / 0")
	testParseHelper(t, math.Inf(-1), "(0 - 1) / 0")
}

func testParseHelper(t *testing.T, answer float64, fStr string, vl ...float64) {
	f, err := Parse(fStr)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	result, err := f(vl...)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	if result != answer {
		t.Logf("Answer is incorrect. %s = %f (%v) Answer: %f\r\n", fStr, result, vl, answer)
		t.Fail()
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
}

func testParseErrorsHelper(t *testing.T, fStr string) {
	_, err := Parse(fStr)
	if err == nil {
		t.Logf("Failed to return error parsing: %s", fStr)
		t.Fail()
	}

	//If error is returned, that is a pass
	// fmt.Printf("Success! %s produced error: %s\r\n", fStr, err.Error())
}

func TestExecutionErrors(t *testing.T) {
	testExecutionErrorsHelper(t, "V0 + V1 + V2", 1, 2)
}

func testExecutionErrorsHelper(t *testing.T, fStr string, vl ...float64) {
	f, err := Parse(fStr)
	if err != nil {
		t.Fail()
	}

	_, err = f(vl...)
	if err == nil {
		t.Fail()
	}
}
