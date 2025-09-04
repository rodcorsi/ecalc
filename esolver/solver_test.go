package esolver

import (
	"math/big"
	"testing"
)

func TestSolve(t *testing.T) {
	assert := func(value string, expected *big.Float, isError bool) {

		x, err := New().Solve(value)
		if (err != nil) != isError {
			t.Errorf("Expected error(%v) value:`%v` result:`%v` Error:%v", isError, value, x, err)
			return
		}

		diff := new(big.Float).Sub(x, expected)
		diff.Abs(diff)
		if diff.Cmp(big.NewFloat(0.0000000001)) > 0 {
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

	assert(".", big.NewFloat(0), false)
	assert("1", big.NewFloat(1), false)
	assert(".1", big.NewFloat(0.1), false)
	assert("1.", big.NewFloat(1), false)
	assert("1.5", big.NewFloat(1.5), false)
	assert("1.5+1", big.NewFloat(2.5), false)
	assert("1.5*2", big.NewFloat(3), false)
	assert("1.5*2+2", big.NewFloat(5), false)
	assert("1+3-2^3/5*3+2*3", big.NewFloat(1+3-8.0/5*3+2*3), false)

	assert(`45d15'25"`, degToDec(`45d15'25"`), false)
	assert(`45d15'`, degToDec(`45d15'`), false)
	assert(`45d25"`, degToDec(`45d25"`), false)
	assert(`15'25"`, degToDec(`15'25"`), false)
	assert(`45d`, degToDec(`45d`), false)
	assert(`45'`, degToDec(`45'`), false)
	assert(`45"`, degToDec(`45"`), false)

	assert(`45d15m25s`, degToDec(`45d15m25s`), false)
	assert(`45d15m`, degToDec(`45d15m`), false)
	assert(`45d25s`, degToDec(`45d25s`), false)
	assert(`15m25s`, degToDec(`15m25s`), false)
	assert(`45m`, degToDec(`45m`), false)
	assert(`45s`, degToDec(`45s`), false)

	assert(`tan45d15'25"`, tan(degToDec(`45d15'25"`)), false)
	assert(`tan45d15'25"+5`, new(big.Float).Add(tan(degToDec(`45d15'25"`)), big.NewFloat(5)), false)
	assert(`tan45d15'25"*5`, new(big.Float).Mul(tan(degToDec(`45d15'25"`)), big.NewFloat(5)), false)
	assert(`5tan45`, new(big.Float).Mul(big.NewFloat(5), tan(degToDec(`45`))), false)
	assert(`tantan30`, tan(tan(big.NewFloat(30))), false)
	assert(`pi3`, new(big.Float).Mul(consts["pi"](), big.NewFloat(3)), false)
	assert(`pipi`, new(big.Float).Mul(consts["pi"](), consts["pi"]()), false)
	assert(`5(5)`, big.NewFloat(25), false)
	assert(`(5)(5)`, big.NewFloat(25), false)
}
