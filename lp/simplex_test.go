package lp_test

import (
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/lp"
	"github.com/stretchr/testify/assert"
)

func TestSimplex(t *testing.T) {
	A := [][]float64{
		{2, 3, 1},
		{4, 1, 2},
		{3, 4, 2},
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
	}
	b := []float64{5, 11, 8, 0, 0, 0}
	c := []float64{5, 4, 3}
	startbasis := []int{3, 4, 5}

	x, y, optimalValue, basis, err := lp.Simplex(c, A, b, startbasis)
	assert.NoError(t, err)

	assert.InDelta(t, 13.0, optimalValue, epsilon)

	expectedX := []float64{2, 0, 1}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedY := []float64{1, 0, 1, 0, 3, 0}
	assert.InDeltaSlice(t, expectedY, y, epsilon)

	expectedBasis := []int{0, 2, 4}
	assert.Equal(t, expectedBasis, basis)
}

func TestPhaseOne(t *testing.T) {
	A := [][]float64{
		{1, 1},
		{4, 1},
		{-3, 1},
	}
	b := []float64{1, 3, -1}

	basis, feasible, err := lp.PhaseOne(A, b, false)
	assert.NoError(t, err)
	assert.True(t, feasible, "Expected feasible solution")

	expectedBasis := []int{2, 4}
	assert.Equal(t, expectedBasis, basis)
}

func NotTestMaximization(t *testing.T) {
	A := [][]float64{
		{0.5, -11.0 / 2.0, -5.0 / 2.0, 9},
		{0.5, -3.0 / 2.0, -0.5, 1},
		{1, 0, 0, 0},
		{-1, 0, 0, 0},
		{0, -1, 0, 0},
		{0, 0, -1, 0},
		{0, 0, 0, -1},
	}
	b := []float64{
		0, 0, 1, 0, 0, 0, 0,
	}

	c := []float64{
		10, -57, -9, -24,
	}

	x, optimalValue, error := lp.Maximize(c, A, b)

	assert.NoError(t, error)

	expectedX := []float64{1, 0, 1, 0}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedOptimalValue := 1.0
	assert.InDelta(t, expectedOptimalValue, optimalValue, epsilon)
}
