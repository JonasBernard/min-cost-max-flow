package lp

import (
	"fmt"
	"math"

	"github.com/JonasBernard/min-cost-max-flow/util"
	"github.com/jedib0t/go-pretty/v6/table"
)

func EllipsoidStep(A [][]float64, a []float64, c []float64, t table.Writer) ([][]float64, []float64) {
	n := len(c)
	nf := float64(n)

	Ac := util.MatMul(A, c)
	cAc := util.DotProduct(c, Ac)
	d := util.ScalarMult(1.0/math.Sqrt(cAc), Ac)

	a_new := util.VectorSub(a, util.ScalarMult(1.0/(nf+1.0), d))

	a_factor := nf * nf / (nf*nf - 1.0)
	A_new := util.ScalarMatrixMult(a_factor, util.MatrixSub(A, util.ScalarMatrixMult(2.0/(nf+1.0), util.ReverseDotProduct(d, d))))

	t.AppendRow(table.Row{
		fmt.Sprintf("%v", Ac),
		fmt.Sprintf("%v", cAc),
		fmt.Sprintf("%v", d),
		fmt.Sprintf("%v", a_new),
		fmt.Sprintf("%v", a_factor),
		fmt.Sprintf("%v", A_new),
	})

	return A_new, a_new
}
