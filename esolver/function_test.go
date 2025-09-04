package esolver

import (
	"math/big"
	"testing"
)

func TestDegToDec(t *testing.T) {
	assert := func(val string, expected *big.Float) {
		result := degToDec(val)
		diff := new(big.Float).Sub(result, expected)
		diff.Abs(diff)
		if diff.Cmp(big.NewFloat(0.0000000001)) > 0 {
			t.Errorf("Error on degToDec expected:`%v` but result was:`%v`\n", expected, val)
		}
	}

	assert("35", dms(big.NewFloat(35), big.NewFloat(0), big.NewFloat(0)))
	assert("35d", dms(big.NewFloat(35), big.NewFloat(0), big.NewFloat(0)))
	assert("35.d", dms(big.NewFloat(35), big.NewFloat(0), big.NewFloat(0)))
	assert("35.5d", dms(big.NewFloat(35.5), big.NewFloat(0), big.NewFloat(0)))
	assert(".5d", dms(big.NewFloat(0.5), big.NewFloat(0), big.NewFloat(0)))
	assert("35d20", dms(big.NewFloat(35), big.NewFloat(20), big.NewFloat(0)))
	assert("35d20'", dms(big.NewFloat(35), big.NewFloat(20), big.NewFloat(0)))
	assert(`35d20'12`, dms(big.NewFloat(35), big.NewFloat(20), big.NewFloat(12)))
	assert(`35d20'12"`, dms(big.NewFloat(35), big.NewFloat(20), big.NewFloat(12)))
	assert(`35d20'12.5"`, dms(big.NewFloat(35), big.NewFloat(20), big.NewFloat(12.5)))
	assert(`20'12.5"`, dms(big.NewFloat(0), big.NewFloat(20), big.NewFloat(12.5)))
	assert(`12.5"`, dms(big.NewFloat(0), big.NewFloat(0), big.NewFloat(12.5)))
}
