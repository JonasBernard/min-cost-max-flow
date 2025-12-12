package lp

import (
	"slices"

	"github.com/JonasBernard/min-cost-max-flow/util"
)

type Basis struct {
	Indices []int
	YValues []float64
}

func (b Basis) injectY(y []float64) {
	for i, bi := range b.Indices {
		b.YValues[bi] = y[i]
	}
}

func (b Basis) ToBasisIndex(index int) int {
	return util.IndexOf(b.Indices, index)
}

// implicitly computes N := {0,...,n-1} \ b.Indices and returns the index of index in N
func (b Basis) ToNonBasisIndex(n int, index int) int {
	i := 0
	for v := range n {
		if slices.Contains(b.Indices, v) {
			continue
		}
		if v == index {
			return i
		}
		i++
	}
	return -1
}

func (b Basis) ToGlobalIndex(basisIndex int) int {
	return b.Indices[basisIndex]
}
