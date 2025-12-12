package lp

import (
	"errors"
	"slices"

	"github.com/JonasBernard/min-cost-max-flow/util"
)

// Returns a feasible start basis on the system Ax <= b, x >= 0, or indicates that the system is infeasible.
func PhaseOne(A [][]float64, b []float64, useDualSimplex bool) (basis []int, feasible bool, err error) {
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

	var optimalValue float64
	var resultbasis []int
	if useDualSimplex {
		err = errors.New("to be implemented correctly")
		return
		_, _, optimalValue, resultbasis, err = DualSimplex(util.Transpose(D), d, util.Neg(c), startbasis)
	} else {
		_, _, optimalValue, resultbasis, err = Simplex(c, D, d, startbasis)
	}

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
