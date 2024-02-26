package network

import (
	"github.com/JonasBernard/min-cost-max-flow/graph"
	"github.com/JonasBernard/min-cost-max-flow/util"
)

/*
Returns the residual graph that is a graph with forward and reverse arcs.
Forward arcs e have capacity C(e) - F(e) where C(e) is the original
capacity of e and F(e) is the flow at e.
Reverse arcs e' have capacity F(-e') where F is the flow and -e' is the
edge that corresponds to e'.

To see where the algorithm comes from, see https://sboyles.github.io/teaching/ce377k/class14.pdf

It is not checked if adding the reverse arcs does not create a multi-graph.
*/
func (n WeigthedNetwork[T]) ResidualGraph(flow map[*graph.WeightedDirectedEdge[T]]float64) (residual WeigthedNetwork[T]) {
	residual = WeigthedNetwork[T]{
		WeigthedDirectedGraph: graph.WeigthedDirectedGraph[T]{
			Vertices: n.Vertices,
			Edges:    make([]*graph.WeightedDirectedEdge[T], 0),
		},
		Source: n.Source,
		Sink:   n.Sink,
	}

	// print("Constructing residual graph...")

	for _, e := range n.Edges {
		flowHere := flow[e]

		// fmt.Printf("Considering edge from %v to %v. Capacity: %v. Flow: %v\n", e.VertexFrom, e.VertexTo, e.Capacity, flowHere)

		if e.Capacity > flowHere {
			// fmt.Printf("adding forward arc for %v and %v of capacity %v\n", e.VertexFrom.Node, e.VertexTo.Node, e.Capacity-flowHere)
			forwardArc := graph.WeightedDirectedEdge[T]{
				VertexFrom:   e.VertexFrom,
				VertexTo:     e.VertexTo,
				Weight:       e.Weight,
				Capacity:     e.Capacity - flowHere,
				IsReverseArc: false,
				OriginalEdge: e,
			}
			residual.Edges = append(residual.Edges, &forwardArc)
		}

		if flowHere > 0 {
			// fmt.Printf("adding reverse arc for %v and %v of capacity %v\n", e.VertexFrom.Node, e.VertexTo.Node, flowHere)
			reverseArc := graph.WeightedDirectedEdge[T]{
				VertexFrom:   e.VertexTo,
				VertexTo:     e.VertexFrom,
				Weight:       -e.Weight,
				Capacity:     flowHere,
				IsReverseArc: true,
				OriginalEdge: e,
			}
			residual.Edges = append(residual.Edges, &reverseArc)
		}

		// fmt.Printf("Resulting graph:\n%v\n", residual)
	}

	// print("Finished constructing residual graph.\n")

	return
}

func (_ WeigthedNetwork[T]) AugmentFlow(flow map[*graph.WeightedDirectedEdge[T]]float64, path graph.WeigthedDirectedGraph[T]) (newFlow map[*graph.WeightedDirectedEdge[T]]float64) {
	newFlow = flow
	pathCapacities := util.MapSlice(path.Edges, func(edge **graph.WeightedDirectedEdge[T]) float64 { return (*edge).Capacity })
	bottleneck := util.MinSlice(pathCapacities)

	// fmt.Printf("Augmenting by bottleneck value  %v\n\n", bottleneck)

	for _, edge := range path.Edges {
		if edge.IsReverseArc {
			newFlow[edge.OriginalEdge] -= bottleneck
		} else {
			newFlow[edge.OriginalEdge] += bottleneck
		}
	}

	// fmt.Printf("New Flow: \n\n")
	// util.PrintMap(newFlow)
	return
}

/*
See algorithm B on page 257 in
https://dl.acm.org/doi/10.1145/321694.321699
*/
func (n WeigthedNetwork[T]) MinCostMaxFlow() (flow map[*graph.WeightedDirectedEdge[T]]float64) {
	flow = make(map[*graph.WeightedDirectedEdge[T]]float64, 0)
	for _, e := range n.Edges {
		flow[e] = 0
	}

	for {
		residual := n.ResidualGraph(flow)
		distances := residual.BellmanFordMoore(residual.Source)

		// fmt.Print("Residual graph:\n")
		// residual.PrintSelfWithFlow(nil)

		// fmt.Printf("Used BellmanFordMoore to find the following distances:\n")
		// util.PrintMap(distances)

		path, err := residual.ShortestPathWithMinHopFromDistances(distances, residual.Source, residual.Sink)
		if err != nil {
			break // no augmenting path found means we are done
		}

		// n.PrintSelfWithFlow(flow)
		// fmt.Printf("Augmenting along: %v\n\n", path.Edges)

		flow = residual.AugmentFlow(flow, *path)
		// n.PrintSelfWithFlow(flow)
	}

	return
}
