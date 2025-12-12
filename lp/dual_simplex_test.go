package lp_test

import (
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/lp"
	"github.com/JonasBernard/min-cost-max-flow/util"
	"github.com/stretchr/testify/assert"
)

func TestDualSimplex(t *testing.T) {
	A := [][]float64{
		{-1, -4, -1, -0},
		{-1, -1, -0, -1},
	}
	b := []float64{-0.5, -1, 0, 0}
	c := []float64{-2, -1}
	startbasis := []int{2, 3}

	x, y, optimalValue, endbasis, err := lp.DualSimplex(util.Transpose(A), b, c, startbasis)
	assert.NoError(t, err)
	assert.InDeltaSlice(t, []float64{1.0 / 6.0, 1.0 / 3.0}, x, epsilon)
	assert.InDeltaSlice(t, []float64{2.0 / 3.0, 1.0 / 3.0, 0, 0}, y, epsilon)
	assert.InDelta(t, -2.0/3.0, optimalValue, epsilon)
	assert.Equal(t, []int{0, 1}, endbasis)
}
