package search_test

import (
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/search"
	"github.com/stretchr/testify/assert"
)

func TestBinarySearch(t *testing.T) {
	A := [][]float64{
		{1, 2},
		{3, 1},
	}
	b := []float64{4, 5}
	c := []float64{1, 1}
	K := 5

	result := search.BinarySearch(A, b, c, K)
	assert.Equal(t, 1, result)
}
