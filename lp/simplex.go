package lp

import (
	"errors"
	"fmt"
	"math"
	"os"
	"slices"

	"github.com/JonasBernard/min-cost-max-flow/util"
	"github.com/jedib0t/go-pretty/v6/table"
)

var ErrUnbounded = errors.New("the linear program is unbounded")

type Basis struct {
	Indices []int
	YValues []float64
}

func (b Basis) injectY(y []float64) {
	for i, bi := range b.Indices {
		b.YValues[bi] = y[i]
	}
}

func (b Basis) ToBasisIndex(index int) int {
	return util.FindValue(b.Indices, index)
}

func (b Basis) ToGlobalIndex(basisIndex int) int {
	return b.Indices[basisIndex]
}

func PhaseOne(A [][]float64, b []float64) (basis []int, feasible bool, err error) {
	pos_rows := util.FindAll(b, func(row float64) bool { return row >= 0 })
	neg_rows := util.FindAll(b, func(row float64) bool { return row < 0 })

	m := len(A)
	n := len(A[0])
	m_plus := len(pos_rows)
	m_minus := len(neg_rows)

	A_plus := util.GetRows(A, pos_rows)
	b_plus := util.GetValues(b, pos_rows)

	A_minus := util.GetRows(A, neg_rows)
	b_minus := util.GetValues(b, neg_rows)

	A_negated := util.NegMatrix(A_minus)
	b_negated := util.Neg(b_minus)

	penalty_A := util.MatMul(util.Transpose(A_minus), util.Ones(m_minus))
	penalty_slack := util.Ones(m_minus)
	c := slices.Concat(penalty_A, penalty_slack)
	c = util.Neg(c)

	d := slices.Concat(b_plus, b_negated, util.Zeros(n+m_minus))

	D := [][]float64{}

	if m_plus > 0 {
		D_1 := util.ConcatColumns(A_plus, util.ZeroMatrix(m_plus, m_minus))
		D = slices.Concat(D, D_1)
	}
	if m_minus > 0 {
		D_2 := util.ConcatColumns(A_negated, util.NegMatrix(util.IdentityMatrix(m_minus)))
		D = slices.Concat(D, D_2)
	}

	D_3 := util.ConcatColumns(util.NegMatrix(util.IdentityMatrix(n)), util.ZeroMatrix(n, m_minus))
	D = slices.Concat(D, D_3)

	if m_minus > 0 {
		D_4 := util.ConcatColumns(util.ZeroMatrix(m_minus, n), util.NegMatrix(util.IdentityMatrix(m_minus)))
		D = slices.Concat(D, D_4)
	}

	startbasis := make([]int, n+m_minus)
	for i := 0; i < n+m_minus; i++ {
		startbasis[i] = m_plus + m_minus + i
	}

	_, _, optimalValue, resultbasis, err := Simplex(c, D, d, startbasis)

	if err != nil {
		return nil, false, err
	}

	if -optimalValue > util.DotProduct(util.Ones(m_minus), b_minus) {
		return nil, false, nil
	}

	basis = make([]int, n)

	j := 0
	for i := range resultbasis {
		if j == n {
			break
		}

		variable := resultbasis[i]
		// exclude slack variables that were constructed before
		if variable >= m+n {
			continue
		}

		// exclude x variables that are not tight, i.e. where the corresponding slack variable is not in the bases
		if variable >= m_plus && variable <= m_plus+m_minus-1 {
			if !slices.Contains(resultbasis, variable+m_minus+n) {
				continue
			}
		}

		basis[j] = variable
		j++
	}

	return basis, true, nil
}

// Simplex method on a system in natural form max c@x s.t. A@x <= b using Bland's pivot rule
func Simplex(c []float64, A [][]float64, b []float64, startbasis []int) (x []float64, y []float64, optimalValue float64, endbasis []int, err error) {
	basis := Basis{
		Indices: make([]int, len(startbasis)),
		YValues: make([]float64, len(A)),
	}
	copy(basis.Indices, startbasis)

	A_B := util.GetRows(A, basis.Indices)
	b_B := util.GetValues(b, basis.Indices)
	x, error := SolveLinearSystem(A_B, b_B)
	if error != nil {
		return nil, nil, 0, nil, error
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Iter", "Sol", "Basis", "Objective", "A_B", "A_B^T", "y_B", "Exit i", "A_B^-1", "w", "Aw", "jSelect", "Enter j", "gamma"})

	iter := 0
	for {
		iter++

		A_B = util.GetRows(A, basis.Indices)
		A_BT := util.Transpose(A_B)

		for i := range basis.YValues {
			basis.YValues[i] = 0
		}

		y_B, error := SolveLinearSystem(A_BT, c)
		if error != nil {
			return nil, nil, 0, basis.Indices, error
		}
		basis.injectY(y_B)

		// This is Blands exiting rule because the first index is found.
		i := util.Find(basis.YValues, func(v float64) bool { return v < 0.0 })

		if i == -1 {
			t.AppendRows([]table.Row{
				{iter, x,
					fmt.Sprint(basis.Indices) + "\n" + fmt.Sprint(util.MapSlice(basis.Indices, func(t *int) int { return *t + 1 })),
					util.DotProduct(c, x), util.PrintMatrix(A_B), util.PrintMatrix(A_BT), util.PrintVector(y_B), i},
			})
			t.AppendFooter(table.Row{"Optimal"})
			t.Render()
			return x, basis.YValues, util.DotProduct(c, x), basis.Indices, nil
		}

		inverted, error := Invert(A_B)
		if error != nil {
			t.AppendFooter(table.Row{"Fatal", "error", "while", "matrix", "inversion"})
			t.Render()
			return nil, nil, 0, nil, error
		}

		w := util.Neg(util.GetColumn(inverted, basis.ToBasisIndex(i)))

		ratioTest := util.MatMul(A, w)

		jSelect := util.FindAll(ratioTest, func(v float64) bool { return v > 0 })
		if len(jSelect) == 0 {
			t.AppendRows([]table.Row{
				{iter, x,
					fmt.Sprint(basis.Indices) + "\n" + fmt.Sprint(util.MapSlice(basis.Indices, func(t *int) int { return *t + 1 })),
					util.DotProduct(c, x), util.PrintMatrix(A_B), util.PrintMatrix(A_BT), util.PrintVector(y_B), i, util.PrintMatrix(inverted), util.PrintVector(w), util.PrintVector(ratioTest), jSelect},
			})
			t.AppendFooter(table.Row{"Unbounded"})
			t.Render()
			return nil, nil, 0, basis.Indices, ErrUnbounded
		}

		// Bland's entering rule
		gamma := math.Inf(1)
		j := -1
		for _, jCandidate := range jSelect {
			if util.Contains(basis.Indices, jCandidate) {
				continue
			}
			ratio := (b[jCandidate] - util.DotProduct(util.GetRow(A, jCandidate), x)) / ratioTest[jCandidate]

			if ratio < gamma {
				gamma = ratio
				j = jCandidate
			}
		}

		gamma_w := util.ScalarMult(gamma, w)

		t.AppendRows([]table.Row{
			{iter,
				x,
				fmt.Sprint(basis.Indices) + "\n" + fmt.Sprint(util.MapSlice(basis.Indices, func(t *int) int { return *t + 1 })),
				util.DotProduct(c, x),
				util.PrintMatrix(A_B),
				util.PrintMatrix(A_BT),
				util.PrintVector(y_B),
				i,
				util.PrintMatrix(inverted),
				util.PrintVector(w),
				util.PrintVector(ratioTest),
				jSelect,
				j,
				gamma},
		})
		t.AppendSeparator()

		// Update basis
		basis.Indices = util.RemoveValue(basis.Indices, i)
		basis.Indices = append(basis.Indices, j)

		slices.Sort(basis.Indices)

		x = util.VectorAdd(x, gamma_w)
	}
}

func Maximize(c []float64, A [][]float64, b []float64) (x []float64, optimalValue float64, err error) {
	basis, feasible, err := PhaseOne(A, b)
	if err != nil {
		return nil, 0, err
	}
	if !feasible {
		return nil, 0, errors.New("the linear program is infeasible")
	}

	x, _, optimalValue, _, err = Simplex(c, A, b, basis)
	if err != nil {
		return nil, 0, err
	}

	return x, optimalValue, nil
}
