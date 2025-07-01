package arithmetic_test

import (
	"testing"

	"github.com/felipewom/go-hexagonal/internal/application/core/arithmetic"
	"github.com/stretchr/testify/assert"
)

func TestArith_Operations(t *testing.T) {
	core := arithmetic.New() // This returns api.Arithmetic, but we are testing the concrete Arith impl.

	testCases := []struct {
		name      string
		op        string // "add", "sub", "mul", "div"
		a         int32
		b         int32
		expected  int32
		expectErr error
	}{
		{"addition positive", "add", 5, 3, 8, nil},
		{"addition negative", "add", -5, 3, -2, nil},
		{"addition zero", "add", 5, 0, 5, nil},

		{"subtraction positive", "sub", 5, 3, 2, nil},
		{"subtraction negative result", "sub", 3, 5, -2, nil},
		{"subtraction zero", "sub", 5, 0, 5, nil},
		{"subtraction from zero", "sub", 0, 5, -5, nil},

		{"multiplication positive", "mul", 5, 3, 15, nil},
		{"multiplication by zero", "mul", 5, 0, 0, nil},
		{"multiplication by negative", "mul", 5, -3, -15, nil},
		{"multiplication two negatives", "mul", -5, -3, 15, nil},

		{"division successful", "div", 6, 3, 2, nil},
		{"division integer truncation", "div", 7, 3, 2, nil},
		{"division by zero", "div", 5, 0, 0, arithmetic.ErrDivisionByZero},
		{"division of zero", "div", 0, 5, 0, nil},
		{"division negative numbers", "div", -6, 3, -2, nil},
		{"division two negatives", "div", -6, -3, 2, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res int32
			var err error

			switch tc.op {
			case "add":
				res, err = core.Addition(tc.a, tc.b)
			case "sub":
				res, err = core.Subtraction(tc.a, tc.b)
			case "mul":
				res, err = core.Multiplication(tc.a, tc.b)
			case "div":
				res, err = core.Division(tc.a, tc.b)
			default:
				t.Fatalf("unknown operation: %s", tc.op)
			}

			if tc.expectErr != nil {
				assert.Error(t, err, "Expected an error for %s", tc.name)
				assert.Equal(t, tc.expectErr, err, "Error message mismatch for %s", tc.name)
				// Check expected value even on error, especially for division by zero
				if tc.name == "division by zero" {
					assert.Equal(t, tc.expected, res, "Value mismatch on error for %s", tc.name)
				}
			} else {
				assert.NoError(t, err, "Did not expect an error for %s", tc.name)
				assert.Equal(t, tc.expected, res, "Result mismatch for %s", tc.name)
			}
		})
	}
}
