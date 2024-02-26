package graph

import "math"

/*
Constructs a map that maps every vertex to the shortest distance from s to it using Bellman-Ford-Moore.
*/
func (g WeigthedDirectedGraph[T]) BellmanFordMoore(s Vertex[T]) (distances map[Vertex[T]]float64) {
	distances = make(map[Vertex[T]]float64)
	for _, v := range g.Vertices {
		if v.Node == s.Node {
			distances[v] = 0
		} else {
			distances[v] = math.Inf(+1)
		}
	}
	n := len(g.Vertices)
	for k := 0; k < n; k++ {
		for _, edge := range g.Edges {
			distances[edge.VertexTo] = math.Min(distances[edge.VertexTo], distances[edge.VertexFrom]+edge.Weight)
		}
	}
	return
}
