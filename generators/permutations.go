package generators

import "math/rand/v2"

func GeneratePermutation(n int) []int {
	result := make([]int, n)
	for i := range result {
		result[i] = i + 1
	}
	for i := range result {
		j := rand.IntN(n - i)
		result[i], result[j+i] = result[j+i], result[i]
	}
	return result
}
