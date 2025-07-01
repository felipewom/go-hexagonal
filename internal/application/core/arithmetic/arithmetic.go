package arithmetic

import (
	"errors" // Import for errors.New
	"github.com/felipewom/go-hexagonal/internal/application/api"
)

// Arith implements the api.Arithmetic interface
type Arith struct {
}

// Compile-time check to ensure Arith implements the api.Arithmetic interface.
var _ api.Arithmetic = (*Arith)(nil)

// ErrDivisionByZero is returned when a division by zero is attempted.
var ErrDivisionByZero = errors.New("division by zero")

// New creates a new Arith instance satisfying the api.Arithmetic interface.
func New() api.Arithmetic {
	return &Arith{}
}

// Addition gets the result of adding parameters a and b
func (arith Arith) Addition(a int32, b int32) (int32, error) {
	return a + b, nil
}

// Subtraction gets the result of subtracting parameters a and b
func (arith Arith) Subtraction(a int32, b int32) (int32, error) {
	return a - b, nil
}

// Multiplication gets the result of multiplying parameters a and b
func (arith Arith) Multiplication(a int32, b int32) (int32, error) {
	return a * b, nil
}

// Division gets the result of dividing parameters a and b
func (arith Arith) Division(a int32, b int32) (int32, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}
