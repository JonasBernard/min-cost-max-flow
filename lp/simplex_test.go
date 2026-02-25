package lp_test

import (
	"math"
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/lp"
	"github.com/JonasBernard/min-cost-max-flow/util"
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

func TestSimplexExample91(t *testing.T) {
	A := [][]float64{
		{1, 0, 1},
		{1, 1, 0},
		{1, 2, 0},
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
	}
	b := []float64{8, 7, 12, 0, 0, 0}
	c := []float64{3, 2, 2}
	startbasis := []int{3, 4, 5}

	x, _, optimalValue, basis, err := lp.Simplex(c, A, b, startbasis)
	assert.NoError(t, err)

	assert.InDelta(t, 28, optimalValue, epsilon)

	expectedX := []float64{2, 5, 6}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedBasis := []int{0, 1, 2}
	assert.Equal(t, expectedBasis, basis)
}

func TestSimplexExample91Dual(t *testing.T) { // broken I don't know
	A := [][]float64{
		{1, 0, 1},
		{1, 1, 0},
		{1, 2, 0},
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
	}
	b := []float64{8, 7, 12, 0, 0, 0}
	c := []float64{3, 2, 2}
	startbasis := []int{3, 4, 5}

	x, _, optimalValue, basis, err := lp.DualSimplex(A, b, c, startbasis)
	assert.NoError(t, err)

	assert.InDelta(t, 28, optimalValue, epsilon)

	expectedX := []float64{2, 5, 6}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedBasis := []int{0, 1, 2}
	assert.Equal(t, expectedBasis, basis)
}

func TestSimplexExample92(t *testing.T) {
	A := [][]float64{
		{110, 205, 160, 160, 420, 260},
		{4, 32, 13, 8, 4, 14},
		{2, 12, 54, 284, 22, 80},
		{-110, -205, -160, -160, -420, -260},
		{-4, -32, -13, -8, -4, -14},
		{-2, -12, -54, -284, -22, -80},
		{-1, 0, 0, 0, 0, 0},
		{0, -1, 0, 0, 0, 0},
		{0, 0, -1, 0, 0, 0},
		{0, 0, 0, -1, 0, 0},
		{0, 0, 0, 0, -1, 0},
		{0, 0, 0, 0, 0, -1},
	}
	b := []float64{2000, 55, 800, -2000, -55, -800, 0, 0, 0, 0, 0, 0}
	c := []float64{-3, -24, -13, -9, -20, -19}
	startbasis := []int{0, 1, 2, 6, 7, 11}

	x, _, optimalValue, basis, err := lp.Simplex(c, A, b, startbasis)
	assert.NoError(t, err)

	assert.InDelta(t, -84.4048, optimalValue, epsilon)

	expectedX := []float64{6.4511, 0, 0, 2.6105, 2.0778, 0}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedBasis := []int{1, 3, 5, 7, 8, 11}
	assert.Equal(t, expectedBasis, basis)
}

func TestSimplexExample93(t *testing.T) {
	A := [][]float64{
		{1, 1, -1, 0},
		{6, 5, 0, -1},
		{0, 0, -1, 0},
		{0, 0, 0, -1},
		{-1, 0, 0, 0},
		{0, -1, 0, 0},
	}
	b := []float64{5, 10, 0, 0, 0, 0}
	c := []float64{36, 30, -3, -4}
	startbasis := []int{2, 3, 4, 5}

	x, _, optimalValue, basis, err := lp.Simplex(c, A, b, startbasis)

	assert.ErrorIs(t, err, lp.ErrUnbounded)

	assert.Nil(t, x)

	assert.Equal(t, math.Inf(1), optimalValue)

	expectedBasis := []int{0, 1, 2, 5}
	assert.Equal(t, expectedBasis, basis)
}

func TestSimplexExample94(t *testing.T) {
	A := [][]float64{
		{0.5, -5.5, -2.5, 9.0},
		{0.5, -1.5, -0.5, 1.0},
		{1.0, 0, 0, 0},
		{-1.0, 0, 0, 0},
		{0, -1.0, 0, 0},
		{0, 0, -1.0, 0},
		{0, 0, 0, -1.0},
	}
	b := []float64{0, 0, 1, 0, 0, 0, 0}
	c := []float64{10, -57, -9, -24}
	startbasis := []int{3, 4, 5, 6}

	x, _, optimalValue, basis, err := lp.Simplex(c, A, b, startbasis)
	assert.NoError(t, err)

	assert.InDelta(t, 1.0, optimalValue, epsilon)

	expectedX := []float64{1, 0, 1, 0}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedBasis := []int{1, 2, 4, 6}
	assert.Equal(t, expectedBasis, basis)
}

func TestSimplexExample141(t *testing.T) {
	A := [][]float64{
		{-1, -3, -1, 0},
		{-1, 1, 0, -1},
	}
	b := []float64{-3, -1}
	c := []float64{-2, -5, 0, 0}
	startbasis := []int{2, 3}

	x, _, optimalValue, basis, err := lp.DualSimplex(util.Transpose(A), c, b, startbasis)
	assert.NoError(t, err)

	assert.InDelta(t, -5.5, optimalValue, epsilon)

	expectedX := []float64{7.0 / 4.0, 1.0 / 4.0}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedBasis := []int{0, 1}
	assert.Equal(t, expectedBasis, basis)
}

func TestSimplexExample142(t *testing.T) {
	A := [][]float64{
		{-1, -3, -1, 0},
		{-1, 1, 0, -1},
		{1, 3, 1, 0},
		{1, -1, 0, 1},
		{-1, 0, 0, 0},
		{0, -1, 0, 0},
		{0, 0, -1, 0},
		{0, 0, 0, -1},
	}
	b := []float64{-3, -1, 3, 1, 0, 0, 0, 0}
	c := []float64{2, 5, 0, 0}
	startbasis := []int{2, 3, 4, 5}

	x, _, optimalValue, basis, err := lp.Simplex(c, A, b, startbasis)
	assert.NoError(t, err)

	assert.InDelta(t, 5.5, optimalValue, epsilon)

	expectedX := []float64{1.5, 0.5, 0, 0}
	assert.InDeltaSlice(t, expectedX, x, epsilon)

	expectedBasis := []int{2, 3, 6, 7}
	assert.Equal(t, expectedBasis, basis)
}
