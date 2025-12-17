package numbers

import (
	"fmt"
	"math"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Reduce(rational float64) (p int, q int) {
	rs := make([]float64, 0)
	ps := make([]int, 2)
	qs := make([]int, 2)

	ps[0] = 0
	ps[1] = 1
	qs[0] = 1
	qs[1] = 0

	rs = append(rs, rational)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"iter", "a_i", "p_i", "q_i", "r_i", "diff"})

	i := 2
	for {
		a := int(math.Floor(rs[i-2]))
		ps = append(ps, a*ps[i-1]+ps[i-2])
		qs = append(qs, a*qs[i-1]+qs[i-2])
		diff := rs[i-2] - float64(a)

		t.AppendRow(table.Row{i - 2, a, ps[i], qs[i], rs[i-2], diff})
		t.AppendSeparator()

		if diff < 1e-10 {
			break
		}
		rs = append(rs, 1/diff)

		if i > 10 {
			fmt.Printf("max iterations reached")
			break
		}
		i++
	}

	t.Render()

	return ps[i], qs[i]
}
