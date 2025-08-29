package main

import (
	"testing"
)

func TestFormatRecurring(t *testing.T) {
	tests := []struct {
		name      string
		n         float64
		precision int
		want      string
	}{
		{"1/3", 1.0 / 3.0, 10, "0.̅3"},
		{"1/6", 1.0 / 6.0, 10, "0.1̅6"},
		{"1/7", 1.0 / 7.0, 14, "0.̅1̅4̅2̅8̅5̅7"},
		{"22/7", 22.0 / 7.0, 14, "3.̅1̅4̅2̅8̅5̅7"},
		{"1/2", 1.0 / 2.0, 10, "0.5"},
		{"-1/3", -1.0 / 3.0, 10, "-0.̅3"},
		{"0", 0, 10, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatRecurring(tt.n, tt.precision); got != tt.want {
				t.Errorf("formatRecurring() = %v, want %v", got, tt.want)
			}
		})
	}
}
