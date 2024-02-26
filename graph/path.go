package graph

import "github.com/JonasBernard/min-cost-max-flow/util"

type Path[T Node] WeigthedDirectedGraph[T]

/*
Checks if the graph is just a path from the given start vertex
aka a tree with root start and only one leaf of depth |V|
*/
func (g WeigthedDirectedGraph[T]) IsPathFrom(start Vertex[T]) bool {
	parents, depths := g.BFS(start, nil)
	for _, v := range g.Vertices {
		_, ok := parents[v]
		if !ok && start.Node != v.Node {
			return false
		}
	}
	depthsAsFloats := util.MapMapValues[Vertex[T], int, float64](depths, func(d int) float64 { return float64(d) })
	return util.MaxMapValue(depthsAsFloats)+1 == float64(len(g.Vertices))
}
