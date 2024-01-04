package matching

import (
	"fmt"
	"math/rand"

	"github.com/JonasBernard/min-cost-max-flow/ctypes"
	"github.com/JonasBernard/min-cost-max-flow/util"
)

type MatchNode[L any, R any] struct {
	Name       string
	IsRight    bool
	IsSource   bool
	IsSink     bool
	LeftValue  L
	RightValue R
}

func (n MatchNode[L, R]) String() string {
	return n.Name
}

type MatchingProblem[L ctypes.Node, R ctypes.Node] struct {
	Lefts  []L
	Rights []R
}

type MatchingEdge[L ctypes.Node, R ctypes.Node] struct {
	Left  L
	Right R
}

func (m MatchingProblem[L, R]) ConstructNetworkFromProblem(weights func(leftNode L, rightNode R) (connect bool, weights float64)) (network ctypes.WeigthedNetwork[MatchNode[L, R]], source ctypes.Vertex[MatchNode[L, R]], sink ctypes.Vertex[MatchNode[L, R]]) {
	leftNodes := util.MapSlice(m.Lefts, func(l *L) ctypes.Vertex[MatchNode[L, R]] {
		return ctypes.V(&MatchNode[L, R]{
			Name:      (*l).String(),
			IsRight:   false,
			LeftValue: *l,
		})
	})

	rightNodes := util.MapSlice(m.Rights, func(r *R) ctypes.Vertex[MatchNode[L, R]] {
		return ctypes.V(&MatchNode[L, R]{
			Name:       (*r).String(),
			IsRight:    true,
			RightValue: *r,
		})
	})

	// the capacity cannot be known beforehand because it depends on what the choices are
	allEdges := make([]*ctypes.WeightedDirectedEdge[MatchNode[L, R]], 0, 4*len(leftNodes)+len(rightNodes))

	source = ctypes.V(&MatchNode[L, R]{
		Name: "S", IsSource: true,
	})

	sink = ctypes.V(&MatchNode[L, R]{
		Name: "T", IsSink: true,
	})

	for _, leftNode := range leftNodes {
		allEdges = append(allEdges, &ctypes.WeightedDirectedEdge[MatchNode[L, R]]{
			VertexFrom: source,
			VertexTo:   leftNode,
			Weight:     1,
			Capacity:   1,
		})

		for _, rightNode := range rightNodes {
			connect, weight := weights(leftNode.Node.LeftValue, rightNode.Node.RightValue)

			if connect {
				newEdge := ctypes.WeightedDirectedEdge[MatchNode[L, R]]{
					VertexFrom: leftNode,
					VertexTo:   rightNode,
					Weight:     weight,
					Capacity:   1,
				}
				allEdges = append(allEdges, &newEdge)
			}
		}
	}

	for _, slot := range rightNodes {
		allEdges = append(allEdges, &ctypes.WeightedDirectedEdge[MatchNode[L, R]]{
			VertexFrom: slot,
			VertexTo:   sink,
			Weight:     1,
			Capacity:   1,
		})
	}

	allVertices := append(leftNodes, rightNodes...)
	allVertices = append(allVertices, source)
	allVertices = append(allVertices, sink)

	network = ctypes.WeigthedNetwork[MatchNode[L, R]]{
		WeigthedDirectedGraph: ctypes.WeigthedDirectedGraph[MatchNode[L, R]]{
			Vertices: allVertices,
			Edges:    allEdges,
		},
		Source: source,
		Sink:   sink,
	}

	return
}

func (m MatchingProblem[L, R]) InterpretNetworkFlow(flow map[*ctypes.WeightedDirectedEdge[MatchNode[L, R]]]float64, source ctypes.Vertex[MatchNode[L, R]], sink ctypes.Vertex[MatchNode[L, R]]) (matching []MatchingEdge[L, R], err error) {
	matchingWeightedEdges := util.FilterMapBoth(flow, func(wde *ctypes.WeightedDirectedEdge[MatchNode[L, R]], f float64) bool {
		if f == 0 {
			return false
		}
		if wde.VertexFrom.Node == source.Node || wde.VertexTo.Node == sink.Node {
			return false
		}
		return true
	})

	matchingEdges := make([]MatchingEdge[L, R], 0, len(matchingWeightedEdges))
	for edge := range matchingWeightedEdges {
		matchingEdges = append(matchingEdges, MatchingEdge[L, R]{
			Left:  edge.VertexFrom.Node.LeftValue,
			Right: edge.VertexTo.Node.RightValue,
		})
	}

	if len(matchingEdges) != len(m.Lefts) {
		err := fmt.Errorf("there is no perfect solution. Computed (one of) the best aproximations (missing %v left and %v right nodes)", len(m.Rights)-len(matchingEdges), len(m.Lefts)-len(matchingEdges))
		return matchingEdges, err
	} else {
		return matchingEdges, nil
	}
}

// Shuffles the order of the given edges iteration times to get different results
func (m MatchingProblem[L, R]) SolveMany(iterations int, weights func(leftNode L, rightNode R) (connect bool, weights float64)) (matchings [][]MatchingEdge[L, R], err error) {
	network, source, sink := m.ConstructNetworkFromProblem(weights)

	for i := 0; i < iterations; i++ {
		edges := network.Edges
		rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
		network.Edges = edges

		flow := network.MinCostMaxFlow()

		matching, errThis := m.InterpretNetworkFlow(flow, source, sink)
		matchings = append(matchings, matching)

		if errThis != nil {
			err = errThis
		}

		fmt.Printf("err: %v\n", err)
	}

	return matchings, err
}

func (m MatchingProblem[L, R]) Solve(weights func(leftNode L, rightNode R) (connect bool, weight float64)) (matching []MatchingEdge[L, R], err error) {
	network, source, sink := m.ConstructNetworkFromProblem(weights)

	flow := network.MinCostMaxFlow()

	return m.InterpretNetworkFlow(flow, source, sink)
}
