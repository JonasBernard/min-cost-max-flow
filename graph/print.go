package graph

import (
	"fmt"
	"strings"

	"github.com/JonasBernard/min-cost-max-flow/util"
)

func (e WeightedDirectedEdge[T]) String() string {
	return fmt.Sprintf("%v %v (wgt %v, cap %v)", e.VertexFrom, e.VertexTo, e.Weight, e.Capacity)
}

func (g WeigthedDirectedGraph[T]) String() string {
	edgeList := strings.Join(util.MapSlice(g.Edges, func(e **WeightedDirectedEdge[T]) string { return fmt.Sprintf("%v", *e) }), "\n")
	return fmt.Sprintf("--- Graph of %v nodes and %v edges ---\n%v\n---\n", len(g.Vertices), len(g.Edges), edgeList)
}

func (g WeigthedDirectedGraph[T]) PrintSelfWithFlow(flow map[*WeightedDirectedEdge[T]]float64) {
	edgeList := strings.Join(util.MapSlice(g.Edges, func(e **WeightedDirectedEdge[T]) string { return fmt.Sprintf("%v -> %v", *e, flow[*e]) }), "\n")
	fmt.Printf("--- Graph of %v nodes and %v edges ---\n%v\n---\n", len(g.Vertices), len(g.Edges), edgeList)
}

func (v Vertex[T]) String() string {
	return fmt.Sprintf("[%v]", (*v.Node).String())
}
