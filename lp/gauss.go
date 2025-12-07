package lp

import (
	"math"

	"github.com/JonasBernard/min-cost-max-flow/util"

	"errors"
)

var ErrInvalidInput = errors.New("invalid input: system dimensions mismatch")
var ErrSingularMatrix = errors.New("no pivot found: singular matrix")

func GaussElimination(A [][]float64, b []float64) (RU [][]float64, bU []float64, err error) {
	N := len(A)
	// Only for square matrices
	if N == 0 || len(A[0]) != N || len(b) != N {
		return nil, nil, ErrInvalidInput
	}

	RU = util.MapSlice(A, func(row *[]float64) []float64 {
		newRow := make([]float64, len(*row))
		copy(newRow, *row)
		return newRow
	})
	bU = util.MapSlice(b, func(v *float64) float64 {
		return *v
	})

	for col := 0; col < N; col++ {
		i_max := col
		for row := col + 1; row < N; row++ {
			if math.Abs(RU[row][col]) > math.Abs(RU[i_max][col]) {
				i_max = row
			}
		}
		if RU[i_max][col] == 0 {
			return nil, nil, ErrSingularMatrix
		}
		RU[col], RU[i_max] = RU[i_max], RU[col]
		bU[col], bU[i_max] = bU[i_max], bU[col]

		for row := col + 1; row < N; row++ {
			ratio := RU[row][col] / RU[col][col]
			for j := col + 1; j < N; j++ {
				RU[row][j] -= ratio * RU[col][j]
			}
			bU[row] -= ratio * bU[col]
			RU[row][col] = 0
		}
	}
	return RU, bU, nil
}

func SolveLinearSystem(A [][]float64, b []float64) (x []float64, err error) {
	RU, b, err := GaussElimination(A, b)

	if err != nil {
		return nil, err
	}

	N := len(RU)
	x = make([]float64, N)

	for i := N - 1; i >= 0; i-- {
		x[i] = b[i]
		for j := i + 1; j < N; j++ {
			x[i] -= RU[i][j] * x[j]
		}
		x[i] /= RU[i][i]
	}
	return x, nil
}

func Invert(A [][]float64) (A_inv [][]float64, err error) {
	N := len(A)
	if N == 0 || len(A[0]) != N {
		return nil, ErrInvalidInput
	}

	A_inv = make([][]float64, N)
	for i := 0; i < N; i++ {
		e_i := make([]float64, N)
		e_i[i] = 1
		x_i, err := SolveLinearSystem(A, e_i)
		if err != nil {
			return nil, err
		}
		A_inv[i] = x_i
	}
	return util.Transpose(A_inv), nil
}
