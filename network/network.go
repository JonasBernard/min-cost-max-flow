package network

import (
	"github.com/JonasBernard/min-cost-max-flow/graph"
)

type WeigthedNetwork[T graph.Node] struct {
	graph.WeigthedDirectedGraph[T]
	Source graph.Vertex[T]
	Sink   graph.Vertex[T]
}
