package ctypes

import (
	"errors"
	"fmt"
	"math"
	"strings"

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

type Vertex[T Node] struct {
	Node    *T
	Visited bool
}

func (v Vertex[T]) String() string {
	return fmt.Sprintf("[%v]", (*v.Node).String())
}

func (v *Vertex[T]) Visit() {
	v.Visited = true
}

type WeightedDirectedEdge[T Node] struct {
	VertexFrom   Vertex[T]
	VertexTo     Vertex[T]
	Weight       float64
	Capacity     float64
	IsReverseArc bool
	OriginalEdge *WeightedDirectedEdge[T]
}

func (e WeightedDirectedEdge[T]) String() string {
	return fmt.Sprintf("%v %v (wgt %v, cap %v)", e.VertexFrom, e.VertexTo, e.Weight, e.Capacity)
	//return fmt.Sprintf("%v--(W:%v|C:%v)-->%v", e.VertexFrom, e.Weight, e.Capacity, e.VertexTo)
}

type WeigthedDirectedGraph[T Node] struct {
	Vertices []Vertex[T]
	Edges    []*WeightedDirectedEdge[T]
}

type Path[T Node] WeigthedDirectedGraph[T]

func (g WeigthedDirectedGraph[T]) String() string {
	edgeList := strings.Join(util.MapSlice(g.Edges, func(e **WeightedDirectedEdge[T]) string { return fmt.Sprintf("%v", *e) }), "\n")
	return fmt.Sprintf("--- Graph of %v nodes and %v edges ---\n%v\n---\n", len(g.Vertices), len(g.Edges), edgeList)
}

func (g WeigthedDirectedGraph[T]) PrintSelfWithFlow(flow map[*WeightedDirectedEdge[T]]float64) {
	edgeList := strings.Join(util.MapSlice(g.Edges, func(e **WeightedDirectedEdge[T]) string { return fmt.Sprintf("%v -> %v", *e, flow[*e]) }), "\n")
	fmt.Printf("--- Graph of %v nodes and %v edges ---\n%v\n---\n", len(g.Vertices), len(g.Edges), edgeList)
}

func (v Vertex[T]) IsInGraph(g WeigthedDirectedGraph[T]) bool {
	for _, vert := range g.Vertices {
		if v.Node == vert.Node {
			return true
		}
	}
	return false
}

func (g WeigthedDirectedGraph[T]) getEdge(from Vertex[T], to Vertex[T]) *WeightedDirectedEdge[T] {
	return util.FilterSlice(g.Edges, func(e *WeightedDirectedEdge[T]) bool {
		return e.VertexFrom.Node == from.Node && e.VertexTo.Node == to.Node
	})[0]
}

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
	// fmt.Printf("Incoming Edges of %v from %v\n are %v\n", v, g.Edges, util.FilterSlice(g.Edges, func(e WeightedDirectedEdge[T]) bool { return e.VertexTo.Node == v.Node }))
	return util.FilterSlice(g.Edges, func(e *WeightedDirectedEdge[T]) bool { return e.VertexTo.Node == v.Node })
}

/*
Performs an iterative depth-first-search on the graph starting at the given root node.
It returns the resulting partent structure and the depths of the vertices.
They both depend on the order in which the edges of the graph are given.
*/
func (p WeigthedDirectedGraph[T]) DFS(root Vertex[T]) (parents map[Vertex[T]]Vertex[T], depths map[Vertex[T]]int) {
	p = p.resetVisited()
	parents = make(map[Vertex[T]]Vertex[T])
	depths = make(map[Vertex[T]]int)

	// fmt.Printf("Now running DFS on root node %v\n", root)

	stack := []Vertex[T]{root}
	depths[root] = 0
	root.Visit()

	for {
		if len(stack) == 0 {
			// 	fmt.Printf("Finished running DFS on root node %v\n", root)
			return
		}

		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, v := range p.NeightboursOf(u) {
			if !v.Visited {
				stack = append(stack, v)
				parents[v] = u
				depths[v] = depths[u] + 1
				v.Visit()
			}
		}
	}
}

/*
Performs a synchronous breath-first-search on the graph starting at the given root node.
If the parameter "find" is non-nil, the algorithm will stop and return early upon reaching that given node find.
It returns the resulting partent structure and the depths of the vertices.
They both depend on the order in which the edges of the graph are given.
*/
func (p WeigthedDirectedGraph[T]) BFS(root Vertex[T], find *Vertex[T]) (parents map[Vertex[T]]Vertex[T], depths map[Vertex[T]]int) {
	p = p.resetVisited()
	parents = make(map[Vertex[T]]Vertex[T])
	depths = make(map[Vertex[T]]int)

	// fmt.Printf("Now running BFS on root node %v\n", root)

	stack := []Vertex[T]{root}
	depths[root] = 0
	root.Visit()

	for {
		if len(stack) == 0 {
			// fmt.Printf("Finished running BFS on root node %v\n", root)
			return
		}

		u := stack[0]
		stack = stack[1:]

		for _, v := range p.NeightboursOf(u) {
			if !v.Visited {
				stack = append(stack, v)
				parents[v] = u
				depths[v] = depths[u] + 1
				v.Visit()

				// Optionally give a node to "find" using BFS that stops the
				// algorithm early when reached
				if find != nil && *find == v {
					return
				}
			}
		}
	}
}

/*
Invokes BFS on the graph starting at node root and stops as soon as it reaches
the to node. Then returns the shortest path from root to "to" in terms of hop distance.
Returns an error if there is no shortest path.
*/
func (p WeigthedDirectedGraph[T]) BFSShortestHopPathTo(root Vertex[T], to Vertex[T]) (*WeigthedDirectedGraph[T], error) {
	parents, _ := p.BFS(root, &to)

	path := WeigthedDirectedGraph[T]{[]Vertex[T]{to}, []*WeightedDirectedEdge[T]{}}
	head := to
	for {
		if head.Node == root.Node {
			return &path, nil
		}

		parent, ok := parents[head]
		if !ok {
			return nil, errors.New("No path from root to vertex")
		}
		path.Vertices = append(path.Vertices, parent)
		path.Edges = append(path.Edges, p.getEdge(parent, head))
		head = parent
	}
}

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
			return nil, errors.New(fmt.Sprintf("No path from %v to %v", s, t))
		}

		//fmt.Printf("possible next edges: %v\n", possibleNextEdges)

		selectedEdge := possibleNextEdges[0]

		// detect negative cycles // not sure if it works
		for _, v := range path.Vertices {
			if selectedEdge.VertexFrom.Node == v.Node {
				return nil, errors.New("Detected a negative cycle")
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

		fmt.Printf("---\ncurrent shortestPathTree: \n%v", shortestPathTree)
		fmt.Printf("current head: %v\n", head)
		fmt.Printf("current heads: %v\n", heads)

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

		fmt.Printf("possible next edges: %v\n", possibleNextEdges)

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
				fmt.Printf("appedning to shortestPathTree.. %v\n", newV)
				shortestPathTree.Vertices = append(shortestPathTree.Vertices, newV)
				depths[newV] = depths[head] + 1
				if edge.VertexFrom.Node != s.Node {
					fmt.Printf("appedning to heads.. %v\n", newV)
					heads = append(heads, newV)
				}
			}
		}
	}

	fmt.Printf("final shortestPathTree: \n%v", shortestPathTree)
	return shortestPathTree.BFSShortestHopPathTo(s, t)
}
