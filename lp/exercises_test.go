package lp_test

import (
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/lp"
	"github.com/stretchr/testify/assert"
)

func TestPhaseOneH81(t *testing.T) {
	A := [][]float64{
		{3, -1},
		{4, 2},
		{-1, -1},
		{-3, 1},
	}
	b := []float64{20, 100, -1, -20}

	basis, feasible, err := lp.PhaseOne(A, b, false)
	assert.NoError(t, err)
	assert.True(t, feasible, "Expected feasible solution")

	expectedBasis := []int{0, 5}
	assert.Equal(t, expectedBasis, basis)

	A = [][]float64{
		{3, -1},
		{4, 2},
		{-1, -1},
		{-3, 1},
		{-1, 0},
		{0, -1},
	}
	b = []float64{20, 100, -1, -20, 0, 0}

	c := []float64{2, 1}

	x, y, optimalValue, endbasis, err := lp.Simplex(c, A, b, basis)

	assert.NoError(t, err)

	expectedX := []float64{14, 22}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedOptimalValue := 50.0
	assert.InDelta(t, expectedOptimalValue, optimalValue, epsilon)

	expectedEndbasis := []int{0, 1}
	assert.Equal(t, expectedEndbasis, endbasis)

	expectedY := []float64{0, 0.5, 0, 0, 0, 0}
	assert.InDeltaSlice(t, expectedY, y, epsilon)
}

func TestPhaseOneH81UsingDualSimplex(t *testing.T) {
	A := [][]float64{
		{3, -1},
		{4, 2},
		{-1, -1},
		{-3, 1},
	}
	b := []float64{20, 100, -1, -20}

	basis, feasible, err := lp.PhaseOne(A, b, false)
	assert.NoError(t, err)
	assert.True(t, feasible, "Expected feasible solution")

	expectedBasis := []int{0, 5}
	assert.Equal(t, expectedBasis, basis)

	A = [][]float64{
		{3, -1},
		{4, 2},
		{-1, -1},
		{-3, 1},
		{-1, 0},
		{0, -1},
	}
	b = []float64{20, 100, -1, -20, 0, 0}

	c := []float64{2, 1}

	x, y, optimalValue, endbasis, err := lp.DualSimplex(A, b, c, basis)

	assert.NoError(t, err)

	expectedX := []float64{14, 22}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedOptimalValue := 50.0
	assert.InDelta(t, expectedOptimalValue, optimalValue, epsilon)

	expectedEndbasis := []int{0, 1}
	assert.Equal(t, expectedEndbasis, endbasis)

	expectedY := []float64{0, 0.5, 0, 0, 0, 0}
	assert.InDeltaSlice(t, expectedY, y, epsilon)
}

func TestDualSimplexOnH82(t *testing.T) {
	A := [][]float64{
		{-1, -1},
		{-3, 1},
		{-1, 0},
		{0, -1},
	}
	b := []float64{-2, -5, 0, 0}
	c := []float64{-3, -1}
	startbasis := []int{2, 3}

	x, y, optimalValue, endbasis, err := lp.DualSimplex(A, b, c, startbasis)
	assert.NoError(t, err)
	assert.InDeltaSlice(t, []float64{1.75, 0.25}, x, epsilon)
	assert.InDeltaSlice(t, []float64{1.5, 0.5, 0, 0}, y, epsilon)
	assert.InDelta(t, -5.5, optimalValue, epsilon)
	assert.Equal(t, []int{0, 1}, endbasis)
}
