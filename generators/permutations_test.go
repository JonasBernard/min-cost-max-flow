package generators_test

import (
	"fmt"
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/generators"
)

func TestGenerateMultiple(t *testing.T) {
	for range 25 {
		TestGenerate(t)
	}
}

func TestGenerate(t *testing.T) {
	n := 10

	perm := generators.GeneratePermutation(n)

	fmt.Printf("Generated permutation: %v\n", perm)

	if len(perm) != n {
		t.Errorf("Expected permutation length %d, got %d", n, len(perm))
	}

	seen := make(map[int]bool)
	for _, v := range perm {
		if v < 0 || v > n {
			t.Errorf("Value %d out of range [0, %d)", v, n)
		}
		if seen[v] {
			t.Errorf("Duplicate value %d in permutation", v)
		}
		seen[v] = true
	}
}
