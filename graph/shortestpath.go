package graph

import (
	"errors"
	"fmt"

	"github.com/JonasBernard/min-cost-max-flow/util"
)

/*
Reconstructs a shortest path from the distances that Bellman-Ford-Moore returned.
Returns an error if no path exists or a negative cycle ocurrs.
*/
func (g WeigthedDirectedGraph[T]) ShortestPathFromDistances(distances map[Vertex[T]]float64, s Vertex[T], t Vertex[T]) (*WeigthedDirectedGraph[T], error) {
	path := WeigthedDirectedGraph[T]{[]Vertex[T]{t}, []*WeightedDirectedEdge[T]{}}
	head := t
	for {
		if head.Node == s.Node {
			return &path, nil
		}

		possibleNextEdges := util.FilterSlice(g.IncomingEdgesOf(head), func(e *WeightedDirectedEdge[T]) bool {
			return e.Weight+distances[e.VertexFrom] == distances[head]
		})

		if len(possibleNextEdges) == 0 {
			return nil, fmt.Errorf("no path from %v to %v", s, t)
		}

		// fmt.Printf("possible next edges: %v\n", possibleNextEdges)

		selectedEdge := possibleNextEdges[0]

		// detect negative cycles // not sure if it works
		for _, v := range path.Vertices {
			if selectedEdge.VertexFrom.Node == v.Node {
				return nil, errors.New("detected a negative cycle")
			}
		}

		path.Vertices = append(path.Vertices, selectedEdge.VertexFrom)
		path.Edges = append(path.Edges, selectedEdge)
		head = selectedEdge.VertexFrom
	}
}

/*
Reconstructs a shortest path from the distances that Bellman-Ford-Moore returned.
Amongst all shotest paths it returns one that attains the smallest hop distance.
Returns an error if no path exists or a negative cycle ocurrs.
*/
func (g WeigthedDirectedGraph[T]) ShortestPathWithMinHopFromDistances(distances map[Vertex[T]]float64, s Vertex[T], t Vertex[T]) (*WeigthedDirectedGraph[T], error) {
	shortestPathTree := WeigthedDirectedGraph[T]{[]Vertex[T]{t}, []*WeightedDirectedEdge[T]{}}
	heads := []Vertex[T]{t}

	depths := make(map[Vertex[T]]int)
	depths[t] = 0

	for {
		if len(heads) == 0 {
			break
		}
		head := heads[0]
		heads = heads[1:]

		// fmt.Printf("---\ncurrent shortestPathTree: \n%v", shortestPathTree)
		// fmt.Printf("current head: %v\n", head)
		// fmt.Printf("current heads: %v\n", heads)

		possibleNextEdges := util.FilterSlice(g.IncomingEdgesOf(head), func(e *WeightedDirectedEdge[T]) bool {
			attainsDistance := e.Weight+distances[e.VertexFrom] == distances[head]
			if !attainsDistance {
				return false
			}
			doesCreateNonPositiveCycle := false
			for _, v := range shortestPathTree.Vertices {
				if e.VertexFrom.Node == v.Node && depths[v] < depths[head]+1 {
					doesCreateNonPositiveCycle = true
					break
				}
			}
			return !doesCreateNonPositiveCycle
		})

		// fmt.Printf("possible next edges: %v\n", possibleNextEdges)

		// if len(possibleNextEdges) == 0 {
		// 	return nil, errors.New(fmt.Sprintf("No path from %v to %v", s, t))
		// }

		for _, edge := range possibleNextEdges {
			newV := edge.VertexFrom

			// check if the shortestPathTree already contains v
			containsV := false
			for _, g := range shortestPathTree.Vertices {
				if g.Node == newV.Node {
					containsV = true
					break
				}
			}

			// the edge should be added in any case...
			shortestPathTree.Edges = append(shortestPathTree.Edges, edge)
			// ..whereas V should only be added if needed
			if !containsV {
				// fmt.Printf("appedning to shortestPathTree.. %v\n", newV)
				shortestPathTree.Vertices = append(shortestPathTree.Vertices, newV)
				depths[newV] = depths[head] + 1
				if edge.VertexFrom.Node != s.Node {
					// fmt.Printf("appedning to heads.. %v\n", newV)
					heads = append(heads, newV)
				}
			}
		}
	}

	// fmt.Printf("final shortestPathTree: \n%v", shortestPathTree)
	return shortestPathTree.BFSShortestHopPathTo(s, t)
}
