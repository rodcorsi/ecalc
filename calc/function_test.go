package calc

import (
	"testing"
)

func TestDegToDec(t *testing.T) {
	assert := func(val string, expected float64) {
		if degToDec(val) != expected {
			t.Errorf("Error on degToDec expected:`%v` but result was:`%v`\n", expected, val)
		}
	}

	assert("35", dms(35, 0, 0))
	assert("35d", dms(35, 0, 0))
	assert("35.d", dms(35, 0, 0))
	assert("35.5d", dms(35, 30, 0))
	assert(".5d", dms(0, 30, 0))
	assert("35d20", dms(35, 20, 0))
	assert("35d20'", dms(35, 20, 0))
	assert(`35d20'12`, dms(35, 20, 12))
	assert(`35d20'12"`, dms(35, 20, 12))
	assert(`35d20'12.5"`, dms(35, 20, 12.5))
	assert(`20'12.5"`, dms(0, 20, 12.5))
	assert(`12.5"`, dms(0, 0, 12.5))
}
