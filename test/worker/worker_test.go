package worker

import (
	"math"
	"testing"

	"github.com/Egorekski/calc_2.0/internal/worker"
)

func TestEvaluateExpression(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		err        bool
		tolerance  float64
	}{
		{"5+3", 8, false, 0.000001},
		{"10/2+3*4", 17, false, 0.000001},
		{"(2+3)*4", 20, false, 0.000001},
		{"sqrt(16)", 4, false, 0.0001},
		{"sin(3.1415)", 0, false, 0.0001},
		{"cos(0)", 1, false, 0.0001},
		{"log(1)", 0, false, 0.0001},
		{"(2+3)*sqrt(9)+cos(0)", 16, false, 0.0001},
		{"invalid+expression", 0, true, 0},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			result, err := worker.EvaluateExpression(test.expression)

			if test.err && err == nil {
				t.Errorf("Expected error for expression: %s", test.expression)
			} else if !test.err && err != nil {
				t.Errorf("Unexpected error for expression: %s, error: %v", test.expression, err)
			}

			if math.Abs(result-test.expected) > test.tolerance {
				t.Errorf("For expression %s, expected %.6f but got %.6f", test.expression, test.expected, result)
			}
		})
	}
}
