package esolver

import (
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	assert := func(text, expected string) {
		scanner := NewScanner(strings.NewReader(text), nil)
		expected = degToDecString(expected)
		token := scanner.Scan()
		if token.Type != NUMBER || token.Value != expected {
			t.Errorf("Expected:{NUMBER %v} result:{%v %v}\n", expected, token.Type, token.Value)
		}
	}

	// test scan NUMBER
	assert("45", "45")
	assert("45.", "45.")
	assert(".45", ".45")
	assert("45.1", "45.1")
	assert("45d", "45d")
	assert("45.1d", "45.1d")
	assert("45d20", "45d20")
	assert("45d20'", "45d20'")
	assert("45d20.5'", "45d20.5'")
	assert("45d20'15", "45d20'15")
	assert(`45d20'15"`, `45d20'15"`)

	assert("45d20m", "45d20'")
	assert("45d20.5m", "45d20.5'")
	assert("45d20m15", "45d20'15")
	assert(`45d20m15s`, `45d20'15"`)
}
