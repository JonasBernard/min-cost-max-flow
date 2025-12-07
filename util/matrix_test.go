package util_test

import (
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/util"
)

func TestGetRowsThenTranspose(t *testing.T) {
	A := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	A = util.GetRows(A, []int{1, 0})

	expected := [][]float64{
		{4, 1},
		{5, 2},
		{6, 3},
	}
	result := util.Transpose(A)
	for i := range expected {
		for j := range expected[i] {
			if result[i][j] != expected[i][j] {
				t.Errorf("Transpose failed at (%d, %d): expected %v, got %v", i, j, expected[i][j], result[i][j])
			}
		}
	}
}
