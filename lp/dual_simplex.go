package lp

import (
	"fmt"
	"math"
	"os"
	"slices"

	"github.com/JonasBernard/min-cost-max-flow/util"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Works on the problem min y@b s.t. y@A = c, y >= 0
func DualSimplex(A [][]float64, b []float64, c []float64, startbasis []int) (x []float64, y []float64, optimalValue float64, endbasis []int, err error) {
	n := len(A[0])
	m := len(A)

	basis := Basis{
		Indices: make([]int, len(startbasis)),
		YValues: make([]float64, m),
	}
	copy(basis.Indices, startbasis)

	for j := 0; j < m; j++ {
		basis.YValues[j] = 0
	}

	A_BT := util.Transpose(util.GetRows(A, basis.Indices))
	y_B, error := SolveLinearSystem(A_BT, c)
	if error != nil {
		return nil, nil, 0, nil, error
	}
	basis.injectY(y_B)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Iter", "Basis", "Nonbasis", "y", "A_B", "x", "A_N", "z_N", "j", "w_B", "i", "gamma", "Dual objective"})

	iter := 0
	for {
		iter++

		A_B := util.GetRows(A, basis.Indices)
		b_B := util.GetValues(b, basis.Indices)
		x, error := SolveLinearSystem(A_B, b_B)
		if error != nil {
			return nil, nil, 0, nil, error
		}

		N := make([]int, m-n)
		k := 0
		for i := range m {
			if !util.Contains(basis.Indices, i) {
				N[k] = i
				k++
			}
		}

		b_N := util.GetValues(b, N)
		A_N := util.GetRows(A, N)
		z_N := util.VectorSub(b_N, util.MatMul(A_N, x))

		j := util.Find(N, func(l int) bool {
			return z_N[basis.ToNonBasisIndex(m, l)] < 0.0
		})

		if iter == 1 {
			j = 1
		}

		if j == -1 {
			t.AppendRow(table.Row{
				iter,
				fmt.Sprint(basis.Indices) + "\n" + fmt.Sprint(util.MapSlice(basis.Indices, func(t *int) int { return *t + 1 })),
				N,
				util.PrintVector(basis.YValues),
				util.PrintMatrix(A_B),
				x,
				util.PrintMatrix(A_N),
				util.PrintVector(z_N),
				"", "", "", "",
				util.DotProduct(basis.YValues, b),
			})
			t.AppendFooter(table.Row{"Optimal"})
			t.Render()
			return x, basis.YValues, util.DotProduct(basis.YValues, b), basis.Indices, nil
		}

		A_jT := util.GetRow(A, j)
		A_BT := util.Transpose(A_B)

		w_B, error := SolveLinearSystem(A_BT, A_jT)
		if error != nil {
			return nil, nil, 0, basis.Indices, error
		}

		iSelect := util.FindAll(w_B, func(w float64) bool { return w > 0.0 })
		if len(iSelect) == 0 {
			t.AppendFooter(table.Row{"Problem", "is", "unbounded", "in", "dual", "formulation"})
			t.Render()
			return nil, nil, math.Inf(-1), nil, ErrUnbounded
		}

		gamma := math.Inf(1)
		i := -1
		for _, iCandidate := range iSelect {
			globalICandidate := basis.ToGlobalIndex(iCandidate)
			if !util.Contains(basis.Indices, globalICandidate) {
				continue
			}

			ratio := basis.YValues[globalICandidate] / w_B[iCandidate]
			if ratio < gamma {
				gamma = ratio
				i = globalICandidate
			}
		}

		t.AppendRows([]table.Row{
			{iter,
				fmt.Sprint(basis.Indices) + "\n" + fmt.Sprint(util.MapSlice(basis.Indices, func(t *int) int { return *t + 1 })),
				N,
				util.PrintVector(basis.YValues),
				util.PrintMatrix(A_B),
				x,
				util.PrintMatrix(A_N),
				util.PrintVector(z_N),
				j,
				util.PrintVector(w_B),
				i,
				gamma,
				util.DotProduct(basis.YValues, b),
			},
		})
		t.AppendSeparator()

		y_B := util.GetValues(basis.YValues, basis.Indices)
		y_B_update := util.VectorSub(y_B, util.ScalarMult(gamma, w_B))
		basis.injectY(y_B_update)

		basis.Indices = util.RemoveValue(basis.Indices, i)
		basis.Indices = append(basis.Indices, j)

		slices.Sort(basis.Indices)

		basis.YValues[j] = gamma

	}
}
