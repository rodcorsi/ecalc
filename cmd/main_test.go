package main

import (
	"math/big"
	"testing"
)

func Test_formatValue(t *testing.T) {
	tests := []struct {
		name  string
		value *big.Float
		want  string
	}{
		{"do not return decimal places", big.NewFloat(1), "1"},
		{"max number", big.NewFloat(1234567890), "1234567890"},
		{"small number", big.NewFloat(0.00000001), "0.00000001"},
		{"float return only significative", big.NewFloat(1.5), "1.5"},
		{"max length 10 chars", big.NewFloat(1 / 3.0), "0.33333333"},
		{"big number max length 10 chars", big.NewFloat(1/3.0 + 500), "500.333333"},
		{"eng notation when big number bigger than 10 chars", big.NewFloat(12345678901), "1.2346e+10"},
		{"eng notation when small number", big.NewFloat(0.000000001), "1.0000e-09"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatValue(tt.value)
			if got != tt.want {
				t.Errorf("formatValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
