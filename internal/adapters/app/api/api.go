package api

import (
	"github.com/felipewom/go-hexagonal/internal/ports"
)

type Adapter struct {
	arith ports.ArithmeticPort
}
