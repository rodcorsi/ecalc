package calc

import (
	"math"
	"testing"
)

func TestSolve(t *testing.T) {
	assert := func(value string, expected float64, isError bool) {

		x, err := Solve(value)
		if (err != nil) != isError {
			t.Errorf("Expected error(%v) value:`%v` result:`%v` Error:%v", isError, value, x, err)
			return
		}

		if x != expected {
			t.Errorf("Value:`%v` expected:`%v` result:`%v`", value, expected, x)
		}
	}
	// assert("/2", 0, false)
	// assert("", -1, true)
	// assert(".2.", -1, true)
	// assert("1.2.", -1, true)
	// assert(".2.1", -1, true)
	// assert("1.2.3", -1, true)
	// assert("-", -1, true)

	assert(".", 0, false)
	assert("1", 1, false)
	assert(".1", .1, false)
	assert("1.", 1, false)
	assert("1.5", 1.5, false)
	assert("1.5+1", 1.5+1, false)
	assert("1.5*2", 1.5*2, false)
	assert("1.5*2+2", 1.5*2+2, false)
	assert("1+3-2^3/5*3+2*3", 1+3-math.Pow(2, 3)/5*3+2*3, false)

	assert(`45d15'25"`, dms(45, 15, 25), false)
	assert(`45d15'`, dms(45, 15, 0), false)
	assert(`45d25"`, dms(45, 0, 25), false)
	assert(`15'25"`, dms(0, 15, 25), false)
	assert(`45d`, dms(45, 0, 0), false)
	assert(`45'`, dms(0, 45, 0), false)
	assert(`45"`, dms(0, 0, 45), false)

	assert(`45d15m25s`, dms(45, 15, 25), false)
	assert(`45d15m`, dms(45, 15, 0), false)
	assert(`45d25s`, dms(45, 0, 25), false)
	assert(`15m25s`, dms(0, 15, 25), false)
	assert(`45m`, dms(0, 45, 0), false)
	assert(`45s`, dms(0, 0, 45), false)

	assert(`tan45d15'25"`, tan(dms(45, 15, 25)), false)
	assert(`tan45d15'25"+5`, tan(dms(45, 15, 25))+5, false)
	assert(`tan45d15'25"*5`, tan(dms(45, 15, 25))*5, false)
	assert(`5tan45`, 5*tan(dms(45, 0, 0)), false)
	assert(`tantan30`, tan(tan(30)), false)
	assert(`pi3`, math.Pi*3, false)
	assert(`pipi`, math.Pi*math.Pi, false)
	assert(`5(5)`, 5*5, false)
	assert(`(5)(5)`, 5*5, false)
}
