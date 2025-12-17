package search

import (
	"math"

	"github.com/JonasBernard/min-cost-max-flow/lp"
	"github.com/JonasBernard/min-cost-max-flow/numbers"
	"github.com/JonasBernard/min-cost-max-flow/util"
)

func BinarySearch(A [][]float64, b []float64, c []float64, K int) int {
	n := len(A[0])

	Q := math.Pow(float64(n*K), float64(4*n*n))
	l := -Q
	u := Q

	for u-l > 1/(4*Q*Q) {
		mid := (l + u) / 2
		if IsFeasible(A, b, c, mid) {
			u = mid
		} else {
			l = mid
		}
	}

	p, q := numbers.Reduce(u)
	return p / q
}

func IsFeasible(A [][]float64, b []float64, c []float64, lambda float64) bool {
	A = append(A, util.Neg(c))
	b = append(b, -lambda)
	_, feasible, _ := lp.PhaseOne(A, b, false)
	return feasible
}
