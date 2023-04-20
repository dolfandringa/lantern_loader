package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func Min[T Number](x, y T) T {
	return T(math.Min(float64(x), float64(y)))
}
