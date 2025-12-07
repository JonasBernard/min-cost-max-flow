package lp_test

import (
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/lp"
)

var epsilon = 1e-9

func TestGaussElimination3D(t *testing.T) {
	A := [][]float64{
		{2, 1, -1},
		{1, -1, -1},
		{2, 2, 1},
	}
	b := []float64{1, 3, 1}

	RU, bModified, err := lp.GaussElimination(A, b)
	if err != nil {
		t.Fatalf("Failed to perform Gauss elimination: %v", err)
	}

	expectedRU := [][]float64{
		{2, 1, -1},
		{0, -3.0 / 2.0, -0.5},
		{0, 0, 5.0 / 3.0},
	}
	expectedB := []float64{1, 5.0 / 2.0, 5.0 / 3.0}

	for i := range expectedRU {
		for j := range expectedRU[i] {
			if RU[i][j] < expectedRU[i][j]-epsilon || RU[i][j] > expectedRU[i][j]+epsilon {
				t.Errorf("Expected RU[%d][%d] = %f, got %f", i, j, expectedRU[i][j], RU[i][j])
			}
		}
		if bModified[i] < expectedB[i]-epsilon || bModified[i] > expectedB[i]+epsilon {
			t.Errorf("Expected b[%d] = %f, got %f", i, expectedB[i], bModified[i])
		}
	}
}

func TestSimpleSystem3D(t *testing.T) {
	A := [][]float64{
		{2, 1, -1},
		{1, -1, -1},
		{2, 2, 1},
	}
	b := []float64{1, 3, 1}

	x, err := lp.SolveLinearSystem(A, b)
	if err != nil {
		t.Fatalf("Failed to solve linear system: %v", err)
	}

	expected := []float64{2, -2, 1}

	for i, v := range expected {
		if x[i] < v-epsilon || x[i] > v+epsilon {
			t.Errorf("Expected x[%d] = %f, got %f", i, v, x[i])
		}
	}
}

func TestSimpleSystem4D(t *testing.T) {
	A := [][]float64{
		{8, 4, 2, 1},
		{12, 4, 1, 0},
		{12, 2, 0, 0},
		{1, 1, 1, 1},
	}
	b := []float64{14, 15, 0, 0}

	x, err := lp.SolveLinearSystem(A, b)
	if err != nil {
		t.Fatalf("Failed to solve linear system: %v", err)
	}

	expected := []float64{-1, 6, 3, -8}

	for i, v := range expected {
		if x[i] < v-epsilon || x[i] > v+epsilon {
			t.Errorf("Expected x[%d] = %f, got %f", i, v, x[i])
		}
	}
}
