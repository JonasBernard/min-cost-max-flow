package graph

import (
	"fmt"

	"github.com/JonasBernard/min-cost-max-flow/util"
)

// For shorter notation
func V[T Node](node *T) Vertex[T] {
	return Vertex[T]{Node: node}
}

func E[T Node](from, to Vertex[T], weight float64, capacity float64) *WeightedDirectedEdge[T] {
	return &WeightedDirectedEdge[T]{VertexFrom: from, VertexTo: to, Weight: weight, Capacity: capacity}
}

type Node interface {
	comparable
	fmt.Stringer
}

/*
	The "Node" attribute of a graph is generalization of the concept of each vertex having an ID.

Two instances of this struct are
considered equal iff their Node attribute points to the same object.

This is useful because most function calls in Go are "call by value",
so we need a way to have multiple instances of the same vertex.

This could for example be realized by using an ID field of type int
where two verteces are considered eqaul iff their ID is the same number.
However, in most use cases you naturally have objects of type some common type
T that are semantically associated with the vertieces in the graph.

In case this, when using this graph package as library, it is more conveniend
to just store a pointer to T in the vertex than to first creating artificial
IDs and after runnning the desired algorithm mapping those IDs back to your actual objects.

If a user of this package whishes to use natural numbers or strings as Nodes (or "IDs"),
theiy can certainly do so by using "int" or "string" as the type T respectively.
*/
type Vertex[T Node] struct {
	Node    *T
	Visited bool
}

type WeightedDirectedEdge[T Node] struct {
	VertexFrom   Vertex[T]
	VertexTo     Vertex[T]
	Weight       float64
	Capacity     float64
	IsReverseArc bool
	OriginalEdge *WeightedDirectedEdge[T]
}

type WeigthedDirectedGraph[T Node] struct {
	Vertices []Vertex[T]
	Edges    []*WeightedDirectedEdge[T]
}

/*
Checks if this vertex is part of a given graph, where technically we check wheather the vertieces
node that acts like an ID is present in some vertex of the given graph.
*/
func (v Vertex[T]) IsInGraph(g WeigthedDirectedGraph[T]) bool {
	for _, vert := range g.Vertices {
		if v.Node == vert.Node {
			return true
		}
	}
	return false
}

/*
	Returns an instance of an edge from two given verteices.

The function assumes such an edge exists in the graph the function is executed on.
This is useful to get the information that is associated with some edge such as the weight without having access to an
instance of the actual edge but only its endpoints.
*/
func (g WeigthedDirectedGraph[T]) getEdge(from Vertex[T], to Vertex[T]) *WeightedDirectedEdge[T] {
	return util.FilterSlice(g.Edges, func(e *WeightedDirectedEdge[T]) bool {
		return e.VertexFrom.Node == from.Node && e.VertexTo.Node == to.Node
	})[0]
}

/* Marks a given vertex visited */
func (v *Vertex[T]) Visit() {
	v.Visited = true
}

/*
Returns an identical graph to the one given, but all vertices are marked unvisited
*/
func (g WeigthedDirectedGraph[T]) resetVisited() WeigthedDirectedGraph[T] {
	for _, v := range g.Vertices {
		v.Visited = false
	}
	return g
}

/*
Returns a slice of all vertices in the graph that are currently marked visited
*/
func (g WeigthedDirectedGraph[T]) GetVisitedVertices() []Vertex[T] {
	return util.FilterSlice(g.Vertices, func(v Vertex[T]) bool { return v.Visited })
}

/*
Returns a slice of all outgoing edges of v in the graph g.
*/
func (g WeigthedDirectedGraph[T]) OutgoingEdgesOf(v Vertex[T]) []*WeightedDirectedEdge[T] {
	return util.FilterSlice(g.Edges, func(e *WeightedDirectedEdge[T]) bool { return e.VertexFrom.Node == v.Node })
}

/*
Returns a slice containing all (outgoing) neighbours (vertices) of v in the graph g.
*/
func (g WeigthedDirectedGraph[T]) NeightboursOf(v Vertex[T]) []Vertex[T] {
	return util.MapSlice(g.OutgoingEdgesOf(v), func(e **WeightedDirectedEdge[T]) Vertex[T] { return (*e).VertexTo })
}

/*
Returns a slice of all incoming edges of v in the graph g.
*/
func (g WeigthedDirectedGraph[T]) IncomingEdgesOf(v Vertex[T]) []*WeightedDirectedEdge[T] {
	return util.FilterSlice(g.Edges, func(e *WeightedDirectedEdge[T]) bool { return e.VertexTo.Node == v.Node })
}
