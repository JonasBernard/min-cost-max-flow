package network_test

import (
	"fmt"
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/graph"
	"github.com/JonasBernard/min-cost-max-flow/network"
	"github.com/JonasBernard/min-cost-max-flow/util"
)

type TestNode struct {
	Name string
}

func (n TestNode) String() string {
	return n.Name
}

func TestTraversal(t *testing.T) {
	fmt.Printf("Testing if DFS and BFS work as expected\n")

	a := graph.V(&TestNode{Name: "A"})
	b := graph.V(&TestNode{Name: "B"})
	c := graph.V(&TestNode{Name: "C"})

	g := graph.WeigthedDirectedGraph[TestNode]{
		Vertices: []graph.Vertex[TestNode]{a, b, c},
		Edges: []*graph.WeightedDirectedEdge[TestNode]{
			{VertexFrom: a, VertexTo: b, Weight: 2},
		},
	}

	print("Executing DFS.\n")
	parents, depths := g.DFS(a)

	print("Parent structure:\n")
	util.PrintMap(parents)
	print("Depths structure:\n")
	util.PrintMap(depths)

	print("Add an edge from B to C.\n")
	g.Edges = append(g.Edges, &graph.WeightedDirectedEdge[TestNode]{VertexFrom: b, VertexTo: c, Weight: 3})

	print("Executing DFS.\n")
	parents, depths = g.DFS(a)

	print("Parent structure:\n")
	util.PrintMap(parents)
	print("Depths structure:\n")
	util.PrintMap(depths)

	print("Executing BFS.\n")
	parents, depths = g.BFS(a, nil)

	print("Parent structure:\n")
	util.PrintMap(parents)
	print("Depths structure:\n")
	util.PrintMap(depths)
}

func TestIsPath(t *testing.T) {
	a := graph.V(&TestNode{Name: "A"})
	b := graph.V(&TestNode{Name: "B"})
	c := graph.V(&TestNode{Name: "C"})

	g := graph.WeigthedDirectedGraph[TestNode]{
		Vertices: []graph.Vertex[TestNode]{a, b, c},
		Edges: []*graph.WeightedDirectedEdge[TestNode]{
			{VertexFrom: a, VertexTo: b, Weight: 2},
		},
	}

	// not connected
	fmt.Printf("Is graph path from a? %v\n", g.IsPathFrom(a)) // false
	fmt.Printf("Is graph path from b? %v\n", g.IsPathFrom(b)) // false
	fmt.Printf("Is graph path from c? %v\n", g.IsPathFrom(c)) // false

	g.Edges = append(g.Edges, &graph.WeightedDirectedEdge[TestNode]{VertexFrom: b, VertexTo: c, Weight: 3})

	// is a path
	fmt.Printf("Is graph path from a? %v\n", g.IsPathFrom(a)) // true
	fmt.Printf("Is graph path from b? %v\n", g.IsPathFrom(b)) // false
	fmt.Printf("Is graph path from c? %v\n", g.IsPathFrom(c)) // false

	d := graph.V(&TestNode{Name: "D"})

	g.Vertices = append(g.Vertices, d)
	g.Edges = append(g.Edges, &graph.WeightedDirectedEdge[TestNode]{VertexFrom: b, VertexTo: c, Weight: 3})

	// has branches
	fmt.Printf("Is graph path from a? %v\n", g.IsPathFrom(a)) // false
	fmt.Printf("Is graph path from b? %v\n", g.IsPathFrom(b)) // false
	fmt.Printf("Is graph path from c? %v\n", g.IsPathFrom(c)) // false
}

func TestShortestPathHop(t *testing.T) {
	a := graph.V(&TestNode{Name: "A"})
	b := graph.V(&TestNode{Name: "B"})
	c := graph.V(&TestNode{Name: "C"})
	d := graph.V(&TestNode{Name: "D"})
	e := graph.V(&TestNode{Name: "E"})
	f := graph.V(&TestNode{Name: "F"})

	g := graph.WeigthedDirectedGraph[TestNode]{
		Vertices: []graph.Vertex[TestNode]{a, b, c, d, e, f},
		Edges: []*graph.WeightedDirectedEdge[TestNode]{
			{VertexFrom: a, VertexTo: b, Weight: 2},
			{VertexFrom: a, VertexTo: c, Weight: 2},
			{VertexFrom: d, VertexTo: f, Weight: 2},
			{VertexFrom: e, VertexTo: f, Weight: 2},
		},
	}

	path, err := g.BFSShortestHopPathTo(a, f)
	if err != nil {
		fmt.Printf("No shortest path from a to f\n") // should not find a path
	} else {
		fmt.Printf("Found shortest path: %v\n", path.Edges)
	}

	print("Adding edge from C to E.\n")
	g.Edges = append(g.Edges, &graph.WeightedDirectedEdge[TestNode]{
		VertexFrom: c, VertexTo: e, Weight: 0,
	})

	path, err = g.BFSShortestHopPathTo(a, f)
	if err != nil {
		fmt.Printf("No shortest path from a to f\n")
	} else {
		fmt.Printf("Found shortest path: %v\n", path.Edges) // should find path a, c, e, f
	}
}

func TestShortestPathWeights(t *testing.T) {
	a := graph.V(&TestNode{Name: "A"})
	b := graph.V(&TestNode{Name: "B"})
	c := graph.V(&TestNode{Name: "C"})
	d := graph.V(&TestNode{Name: "D"})
	e := graph.V(&TestNode{Name: "E"})
	f := graph.V(&TestNode{Name: "F"})

	// https://upload.wikimedia.org/wikipedia/commons/thumb/3/3b/Shortest_path_with_direct_weights.svg/1200px-Shortest_path_with_direct_weights.svg.png
	graph := graph.WeigthedDirectedGraph[TestNode]{
		Vertices: []graph.Vertex[TestNode]{a, b, c, d, e, f},
		Edges: []*graph.WeightedDirectedEdge[TestNode]{
			// {VertexFrom: a, VertexTo: b, Weight: 4},
			{VertexFrom: a, VertexTo: c, Weight: 2},
			{VertexFrom: c, VertexTo: e, Weight: 3},
			{VertexFrom: b, VertexTo: d, Weight: 10},
			{VertexFrom: b, VertexTo: c, Weight: 5},
			// {VertexFrom: e, VertexTo: d, Weight: 4},
			{VertexFrom: d, VertexTo: f, Weight: 11},
		},
	}

	distances := graph.BellmanFordMoore(a)
	path, err := graph.ShortestPathFromDistances(distances, a, f)
	if err != nil {
		fmt.Printf("No shortest path from a to f\n")
	} else {
		fmt.Printf("Found shortest path: %v\n", path.Edges) // should find path a, c, e, f
	}
}

func TestNegativeCycle(t *testing.T) {
	a := graph.V(&TestNode{Name: "A"})
	b := graph.V(&TestNode{Name: "B"})
	c := graph.V(&TestNode{Name: "C"})
	d := graph.V(&TestNode{Name: "D"})

	graph := graph.WeigthedDirectedGraph[TestNode]{
		Vertices: []graph.Vertex[TestNode]{a, b, c, d},
		Edges: []*graph.WeightedDirectedEdge[TestNode]{
			{VertexFrom: a, VertexTo: b, Weight: 2},
			{VertexFrom: b, VertexTo: c, Weight: -1},
			{VertexFrom: c, VertexTo: b, Weight: -10},
			{VertexFrom: c, VertexTo: d, Weight: 5},
		},
	}

	distances := graph.BellmanFordMoore(a)
	path, err := graph.ShortestPathFromDistances(distances, a, d)
	if err != nil {
		fmt.Printf("No path found: %v\n", err)
	} else {
		fmt.Printf("Found shortest path: %v\n", path.Edges) // should find path a, c, e, f
	}
}

func TestResidualGraph(test *testing.T) {
	a := graph.V(&TestNode{Name: "A"})
	b := graph.V(&TestNode{Name: "B"})

	d := graph.V(&TestNode{Name: "C"})
	c := graph.V(&TestNode{Name: "D"})

	s := graph.V(&TestNode{Name: "S"})
	t := graph.V(&TestNode{Name: "T"})

	e1 := graph.E(s, a, 1, 5)
	e2 := graph.E(s, b, 1, 5)
	e3 := graph.E(a, c, 1, 3)
	e4 := graph.E(a, d, 1, 2)
	e5 := graph.E(b, c, 1, 1)
	e6 := graph.E(b, d, 1, 7)
	e7 := graph.E(c, t, 1, 2)
	e8 := graph.E(d, t, 1, 2)

	network := network.WeigthedNetwork[TestNode]{
		WeigthedDirectedGraph: graph.WeigthedDirectedGraph[TestNode]{
			Vertices: []graph.Vertex[TestNode]{a, b, c, d, s, t},
			Edges:    []*graph.WeightedDirectedEdge[TestNode]{e1, e2, e3, e4, e5, e6, e7, e8},
		},
		Source: s,
		Sink:   t,
	}

	flow := make(map[*graph.WeightedDirectedEdge[TestNode]]float64)
	flow[e1] = 3
	flow[e3] = 1
	flow[e4] = 2
	flow[e7] = 1
	flow[e8] = 2

	util.PrintMap(flow)

	network.PrintSelfWithFlow(flow)

	residual := network.ResidualGraph(flow)

	residual.PrintSelfWithFlow(flow)
}

func TestMinCostFlow(test *testing.T) {
	a := graph.V(&TestNode{Name: "A"})
	b := graph.V(&TestNode{Name: "B"})

	s := graph.V(&TestNode{Name: "S"})
	t := graph.V(&TestNode{Name: "T"})

	network := network.WeigthedNetwork[TestNode]{
		WeigthedDirectedGraph: graph.WeigthedDirectedGraph[TestNode]{
			Vertices: []graph.Vertex[TestNode]{a, b, s, t},
			Edges: []*graph.WeightedDirectedEdge[TestNode]{
				// short notation for edges
				graph.E(s, a, 1, 5),
				graph.E(s, b, 1, 3),

				graph.E(b, a, 1, 1),

				graph.E(a, t, 1, 4),
				graph.E(b, t, 1, 4),
			},
		},
		Source: s,
		Sink:   t,
	}

	flow := network.MinCostMaxFlow()

	print("Computed maximal flow:\n")
	util.PrintMap(flow)
}

func TestMinCostFlow2(test *testing.T) {
	a := graph.V(&TestNode{Name: "2"})
	b := graph.V(&TestNode{Name: "3"})

	s := graph.V(&TestNode{Name: "1"})
	t := graph.V(&TestNode{Name: "4"})

	network := network.WeigthedNetwork[TestNode]{
		WeigthedDirectedGraph: graph.WeigthedDirectedGraph[TestNode]{
			Vertices: []graph.Vertex[TestNode]{a, b, s, t},
			Edges: []*graph.WeightedDirectedEdge[TestNode]{
				// short notation for edges
				graph.E(s, a, 1, 1),
				graph.E(s, b, 5, 3),

				graph.E(a, b, 1, 2),

				graph.E(a, t, 4, 1),
				graph.E(b, t, 2, 3),
			},
		},
		Source: s,
		Sink:   t,
	}

	flow := network.MinCostMaxFlow()

	print("Computed maximal flow:\n")
	util.PrintMap(flow)
}
