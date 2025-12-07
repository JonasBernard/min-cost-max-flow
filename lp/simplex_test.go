package lp_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/lp"
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
	if err != nil {
		t.Fatalf("Simplex returned an error: %v", err)
	}

	expectedOptimalValue := 13.0
	if math.Abs(optimalValue-expectedOptimalValue) > epsilon {
		t.Errorf("Expected optimal value %v, got %v", expectedOptimalValue, optimalValue)
	}

	expectedX := []float64{2, 0, 1}
	for i := range x {
		if math.Abs(x[i]-expectedX[i]) > epsilon {
			t.Errorf("Expected x[%d] = %v, got %v", i, expectedX[i], x[i])
		}
	}

	expectedY := []float64{1, 0, 1, 0, 3, 0}
	for i := range y {
		if math.Abs(y[i]-expectedY[i]) > epsilon {
			t.Errorf("Expected y[%d] = %v, got %v", i, expectedY[i], y[i])
		}
	}

	expectedBasis := []int{0, 2, 4}
	for i := range basis {
		if expectedBasis[i] != basis[i] {
			t.Errorf("Expected basis[%d] = %d, got %d", i, expectedBasis[i], basis[i])
		}
	}
}

func TestPhaseOne(t *testing.T) {
	A := [][]float64{
		{1, 1},
		{4, 1},
		{-3, 1},
	}
	b := []float64{1, 3, -1}

	basis, feasible, err := lp.PhaseOne(A, b)
	if err != nil {
		t.Fatalf("PhaseOne returned an error: %v", err)
	}

	if !feasible {
		t.Fatalf("Expected feasible solution, got infeasible")
	}

	fmt.Printf("%v", basis)

	expectedBasis := []int{2, 4}
	for i := range basis {
		if basis[i] != expectedBasis[i] {
			t.Errorf("Expected basis[%d] = %v, got %v", i, expectedBasis[i], basis[i])
		}
	}
}

func TestMaximization(t *testing.T) {
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

	if error != nil {
		t.Errorf("Maximizazion failed: %v", error)
	}

	expectedX := []float64{
		1, 0, 1, 0,
	}

	for i := range expectedX {
		if math.Abs(x[i]-expectedX[i]) > epsilon {
			t.Errorf("Expected x[%d] = %v, got %v", i, expectedX[i], x[i])
		}
	}

	expectedOptimalValue := 1.0

	if math.Abs(expectedOptimalValue-optimalValue) > epsilon {
		t.Errorf("Expected optimalValue = %v, got %v", expectedOptimalValue, optimalValue)
	}
}
