package ecalc

import (
	"math/big"
	"testing"
)

func TestFormatRecurring(t *testing.T) {
	tests := []struct {
		name      string
		n         *big.Float
		precision int
		want      string
	}{
		{"1/3", new(big.Float).Quo(big.NewFloat(1.0), big.NewFloat(3.0)), 10, "0.̅3"},
		{"1/6", new(big.Float).Quo(big.NewFloat(1.0), big.NewFloat(6.0)), 10, "0.1̅6"},
		{"1/7", new(big.Float).Quo(big.NewFloat(1.0), big.NewFloat(7.0)), 14, "0.̅1̅4̅2̅8̅5̅7"},
		{"22/7", new(big.Float).Quo(big.NewFloat(22.0), big.NewFloat(7.0)), 14, "3.̅1̅4̅2̅8̅5̅7"},
		{"1/2", new(big.Float).Quo(big.NewFloat(1.0), big.NewFloat(2.0)), 10, "0.5"},
		{"-1/3", new(big.Float).Quo(big.NewFloat(-1.0), big.NewFloat(3.0)), 10, "-0.̅3"},
		{"0.332", big.NewFloat(0.332), 10, "0.332"},
		{"0.9999 the precision is equal in all 9 case", big.NewFloat(0.9999), 4, "0.9999"},
		{"0.9999 the precision is small in all 9 case", big.NewFloat(0.9999), 3, "1"},
		{"0", big.NewFloat(0), 10, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatRecurring(tt.n, tt.precision); got != tt.want {
				t.Errorf("formatRecurring() = %v, want %v", got, tt.want)
			}
		})
	}
}
