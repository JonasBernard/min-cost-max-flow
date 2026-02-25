package util

import (
	"fmt"
	"slices"
)

func Transpose(A [][]float64) [][]float64 {
	if len(A) == 0 {
		return [][]float64{}
	}
	M := len(A)
	N := len(A[0])
	AT := make([][]float64, N)
	for i := range N {
		AT[i] = make([]float64, M)
		for j := range M {
			AT[i][j] = A[j][i]
		}
	}
	return AT
}

func GetRows(A [][]float64, rowIndices []int) [][]float64 {
	rows := make([][]float64, len(rowIndices))
	for i, rowIndex := range rowIndices {
		rows[i] = A[rowIndex]
	}
	return rows
}

func GetValues(a []float64, indices []int) []float64 {
	rows := make([]float64, len(indices))
	for i, rowIndex := range indices {
		rows[i] = a[rowIndex]
	}
	return rows
}

func GetRow(A [][]float64, rowIndex int) []float64 {
	return A[rowIndex]
}

func GetColumn(A [][]float64, rowIndex int) []float64 {
	col := make([]float64, len(A))
	for i := range A {
		col[i] = A[i][rowIndex]
	}
	return col
}

func Neg(v []float64) []float64 {
	return MapSlice(v, func(x *float64) float64 {
		return -*x
	})
}

func NegMatrix(A [][]float64) [][]float64 {
	return MapSlice(A, func(t *[]float64) []float64 { return Neg(*t) })
}

func MatMul(A [][]float64, v []float64) []float64 {
	M := len(A)
	N := len(A[0])
	if len(v) != N {
		panic(fmt.Sprintf("Matrix and vector dimensions do not match: %dx%d and %d", M, N, len(v)))
	}
	result := make([]float64, M)
	for i := range M {
		sum := 0.0
		for j := range N {
			sum += A[i][j] * v[j]
		}
		result[i] = sum
	}
	return result
}

func DotProduct(a []float64, b []float64) float64 {
	if len(a) != len(b) {
		panic("Vectors must be of the same length for dot product: got lengths " + fmt.Sprint(len(a)) + " and " + fmt.Sprint(len(b)))
	}
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

func ReverseDotProduct(a []float64, b []float64) [][]float64 {
	if len(a) != len(b) {
		panic("Vectors must be of the same length for reverse dot product: got lengths " + fmt.Sprint(len(a)) + " and " + fmt.Sprint(len(b)))
	}
	n := len(a)
	result := make([][]float64, n)
	for i := range n {
		result[i] = make([]float64, n)
		for j := range n {
			result[i][j] = a[i] * b[j]
		}
	}
	return result
}

func ScalarMult(scalar float64, v []float64) []float64 {
	return MapSlice(v, func(x *float64) float64 {
		return scalar * *x
	})
}

func ScalarMatrixMult(scalar float64, A [][]float64) [][]float64 {
	return MapSlice(A, func(row *[]float64) []float64 {
		return ScalarMult(scalar, *row)
	})
}

func VectorAdd(a []float64, b []float64) []float64 {
	if len(a) != len(b) {
		panic("Vectors must be of the same length for addition")
	}
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result
}

func VectorSub(a []float64, b []float64) []float64 {
	if len(a) != len(b) {
		panic("Vectors must be of the same length for subtraction")
	}
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] - b[i]
	}
	return result
}

func MatrixSub(A [][]float64, B [][]float64) [][]float64 {
	if len(A) != len(B) || len(A[0]) != len(B[0]) {
		panic("Matrices must be of the same dimensions for subtraction")
	}
	M := len(A)
	N := len(A[0])
	result := make([][]float64, M)
	for i := range M {
		result[i] = make([]float64, N)
		for j := range N {
			result[i][j] = A[i][j] - B[i][j]
		}
	}
	return result
}

func Ones(n int) []float64 {
	ones := make([]float64, n)
	for i := range n {
		ones[i] = 1.0
	}
	return ones
}

func Zeros(n int) []float64 {
	zeros := make([]float64, n)
	for i := range n {
		zeros[i] = 0.0
	}
	return zeros
}

func ZeroMatrix(m int, n int) [][]float64 {
	A := make([][]float64, m)
	for i := range A {
		A[i] = make([]float64, n)
	}
	return A
}

func IdentityMatrix(n int) [][]float64 {
	I := ZeroMatrix(n, n)
	for i := range n {
		I[i][i] = 1.0
	}
	return I
}

func ConcatColumns(A [][]float64, B [][]float64) [][]float64 {
	if len(A) == 0 {
		return B
	}
	if len(B) == 0 {
		return A
	}
	if len(A) != len(B) {
		panic("Matrices must have the same number of rows to concatenate")
	}
	M := len(A)
	C := make([][]float64, M)
	for i := range M {
		C[i] = slices.Concat(A[i], B[i])
	}
	return C
}

func PrintMatrix(A [][]float64) string {
	result := ""
	for _, row := range A {
		result += fmt.Sprintf("%.2f\n", row)
	}
	return result
}

func PrintVector(v []float64) string {
	result := ""
	for _, val := range v {
		result += fmt.Sprintf("%.2f\n", val)
	}
	return result
}
