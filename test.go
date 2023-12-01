package main

import (
	"fmt"

	"github.com/JonasBernard/min-cost-max-flow/ctypes"
	"github.com/JonasBernard/min-cost-max-flow/util"
)

type TestNode struct {
	Name string
}

func (n TestNode) String() string {
	return n.Name
}

func testTraversal() {
	fmt.Printf("Testing if DFS and BFS work as expected\n")

	a := ctypes.V(&TestNode{Name: "A"})
	b := ctypes.V(&TestNode{Name: "B"})
	c := ctypes.V(&TestNode{Name: "C"})

	graph := ctypes.WeigthedDirectedGraph[TestNode]{
		Vertices: []ctypes.Vertex[TestNode]{a, b, c},
		Edges: []*ctypes.WeightedDirectedEdge[TestNode]{
			{VertexFrom: a, VertexTo: b, Weight: 2},
		},
	}

	print("Executing DFS.\n")
	parents, depths := graph.DFS(a)

	print("Parent structure:\n")
	util.PrintMap(parents)
	print("Depths structure:\n")
	util.PrintMap(depths)

	print("Add an edge from B to C.\n")
	graph.Edges = append(graph.Edges, &ctypes.WeightedDirectedEdge[TestNode]{VertexFrom: b, VertexTo: c, Weight: 3})

	print("Executing DFS.\n")
	parents, depths = graph.DFS(a)

	print("Parent structure:\n")
	util.PrintMap(parents)
	print("Depths structure:\n")
	util.PrintMap(depths)

	print("Executing BFS.\n")
	parents, depths = graph.BFS(a, nil)

	print("Parent structure:\n")
	util.PrintMap(parents)
	print("Depths structure:\n")
	util.PrintMap(depths)
}

func testIsPath() {
	a := ctypes.V(&TestNode{Name: "A"})
	b := ctypes.V(&TestNode{Name: "B"})
	c := ctypes.V(&TestNode{Name: "C"})

	graph := ctypes.WeigthedDirectedGraph[TestNode]{
		Vertices: []ctypes.Vertex[TestNode]{a, b, c},
		Edges: []*ctypes.WeightedDirectedEdge[TestNode]{
			{VertexFrom: a, VertexTo: b, Weight: 2},
		},
	}

	// not connected
	fmt.Printf("Is graph path from a? %v\n", graph.IsPathFrom(a)) // false
	fmt.Printf("Is graph path from b? %v\n", graph.IsPathFrom(b)) // false
	fmt.Printf("Is graph path from c? %v\n", graph.IsPathFrom(c)) // false

	graph.Edges = append(graph.Edges, &ctypes.WeightedDirectedEdge[TestNode]{VertexFrom: b, VertexTo: c, Weight: 3})

	// is a path
	fmt.Printf("Is graph path from a? %v\n", graph.IsPathFrom(a)) // true
	fmt.Printf("Is graph path from b? %v\n", graph.IsPathFrom(b)) // false
	fmt.Printf("Is graph path from c? %v\n", graph.IsPathFrom(c)) // false

	d := ctypes.V(&TestNode{Name: "D"})

	graph.Vertices = append(graph.Vertices, d)
	graph.Edges = append(graph.Edges, &ctypes.WeightedDirectedEdge[TestNode]{VertexFrom: b, VertexTo: c, Weight: 3})

	// has branches
	fmt.Printf("Is graph path from a? %v\n", graph.IsPathFrom(a)) // false
	fmt.Printf("Is graph path from b? %v\n", graph.IsPathFrom(b)) // false
	fmt.Printf("Is graph path from c? %v\n", graph.IsPathFrom(c)) // false
}

func testShortestPathHop() {
	a := ctypes.V(&TestNode{Name: "A"})
	b := ctypes.V(&TestNode{Name: "B"})
	c := ctypes.V(&TestNode{Name: "C"})
	d := ctypes.V(&TestNode{Name: "D"})
	e := ctypes.V(&TestNode{Name: "E"})
	f := ctypes.V(&TestNode{Name: "F"})

	graph := ctypes.WeigthedDirectedGraph[TestNode]{
		Vertices: []ctypes.Vertex[TestNode]{a, b, c, d, e, f},
		Edges: []*ctypes.WeightedDirectedEdge[TestNode]{
			{VertexFrom: a, VertexTo: b, Weight: 2},
			{VertexFrom: a, VertexTo: c, Weight: 2},
			{VertexFrom: d, VertexTo: f, Weight: 2},
			{VertexFrom: e, VertexTo: f, Weight: 2},
		},
	}

	path, err := graph.BFSShortestHopPathTo(a, f)
	if err != nil {
		fmt.Printf("No shortest path from a to f\n") // should not find a path
	} else {
		fmt.Printf("Found shortest path: %v\n", path.Edges)
	}

	print("Adding edge from C to E.\n")
	graph.Edges = append(graph.Edges, &ctypes.WeightedDirectedEdge[TestNode]{
		VertexFrom: c, VertexTo: e, Weight: 0,
	})

	path, err = graph.BFSShortestHopPathTo(a, f)
	if err != nil {
		fmt.Printf("No shortest path from a to f\n")
	} else {
		fmt.Printf("Found shortest path: %v\n", path.Edges) // should find path a, c, e, f
	}
}

func testShortestPathWeights() {
	a := ctypes.V(&TestNode{Name: "A"})
	b := ctypes.V(&TestNode{Name: "B"})
	c := ctypes.V(&TestNode{Name: "C"})
	d := ctypes.V(&TestNode{Name: "D"})
	e := ctypes.V(&TestNode{Name: "E"})
	f := ctypes.V(&TestNode{Name: "F"})

	// https://upload.wikimedia.org/wikipedia/commons/thumb/3/3b/Shortest_path_with_direct_weights.svg/1200px-Shortest_path_with_direct_weights.svg.png
	graph := ctypes.WeigthedDirectedGraph[TestNode]{
		Vertices: []ctypes.Vertex[TestNode]{a, b, c, d, e, f},
		Edges: []*ctypes.WeightedDirectedEdge[TestNode]{
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

func testNegativeCycle() {
	a := ctypes.V(&TestNode{Name: "A"})
	b := ctypes.V(&TestNode{Name: "B"})
	c := ctypes.V(&TestNode{Name: "C"})
	d := ctypes.V(&TestNode{Name: "D"})

	graph := ctypes.WeigthedDirectedGraph[TestNode]{
		Vertices: []ctypes.Vertex[TestNode]{a, b, c, d},
		Edges: []*ctypes.WeightedDirectedEdge[TestNode]{
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

func testResidualGraph() {
	a := ctypes.V(&TestNode{Name: "A"})
	b := ctypes.V(&TestNode{Name: "B"})

	d := ctypes.V(&TestNode{Name: "C"})
	c := ctypes.V(&TestNode{Name: "D"})

	s := ctypes.V(&TestNode{Name: "S"})
	t := ctypes.V(&TestNode{Name: "T"})

	e1 := ctypes.E(s, a, 1, 5)
	e2 := ctypes.E(s, b, 1, 5)
	e3 := ctypes.E(a, c, 1, 3)
	e4 := ctypes.E(a, d, 1, 2)
	e5 := ctypes.E(b, c, 1, 1)
	e6 := ctypes.E(b, d, 1, 7)
	e7 := ctypes.E(c, t, 1, 2)
	e8 := ctypes.E(d, t, 1, 2)

	network := ctypes.WeigthedNetwork[TestNode]{
		WeigthedDirectedGraph: ctypes.WeigthedDirectedGraph[TestNode]{
			Vertices: []ctypes.Vertex[TestNode]{a, b, c, d, s, t},
			Edges:    []*ctypes.WeightedDirectedEdge[TestNode]{e1, e2, e3, e4, e5, e6, e7, e8},
		},
		Source: s,
		Sink:   t,
	}

	flow := make(map[*ctypes.WeightedDirectedEdge[TestNode]]float64)
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

func testMinCostFlow() {
	a := ctypes.V(&TestNode{Name: "A"})
	b := ctypes.V(&TestNode{Name: "B"})

	s := ctypes.V(&TestNode{Name: "S"})
	t := ctypes.V(&TestNode{Name: "T"})

	network := ctypes.WeigthedNetwork[TestNode]{
		WeigthedDirectedGraph: ctypes.WeigthedDirectedGraph[TestNode]{
			Vertices: []ctypes.Vertex[TestNode]{a, b, s, t},
			Edges: []*ctypes.WeightedDirectedEdge[TestNode]{
				// short notation for edges
				ctypes.E(s, a, 1, 5),
				ctypes.E(s, b, 1, 3),

				ctypes.E(b, a, 1, 1),

				ctypes.E(a, t, 1, 4),
				ctypes.E(b, t, 1, 4),
			},
		},
		Source: s,
		Sink:   t,
	}

	flow := network.MinCostMaxFlow()

	print("Computed maximal flow:\n")
	util.PrintMap(flow)
}

func testMinCostFlow2() {
	a := ctypes.V(&TestNode{Name: "2"})
	b := ctypes.V(&TestNode{Name: "3"})

	s := ctypes.V(&TestNode{Name: "1"})
	t := ctypes.V(&TestNode{Name: "4"})

	network := ctypes.WeigthedNetwork[TestNode]{
		WeigthedDirectedGraph: ctypes.WeigthedDirectedGraph[TestNode]{
			Vertices: []ctypes.Vertex[TestNode]{a, b, s, t},
			Edges: []*ctypes.WeightedDirectedEdge[TestNode]{
				// short notation for edges
				ctypes.E(s, a, 1, 1),
				ctypes.E(s, b, 5, 3),

				ctypes.E(a, b, 1, 2),

				ctypes.E(a, t, 4, 1),
				ctypes.E(b, t, 2, 3),
			},
		},
		Source: s,
		Sink:   t,
	}

	flow := network.MinCostMaxFlow()

	print("Computed maximal flow:\n")
	util.PrintMap(flow)
}

func _main() {
	// testTraversal()
	// testIsPath()
	// testShortestPathHop()
	// testShortestPathWeights()
	// testNegativeCycle()
	// testResidualGraph()
	// testMinCostFlow()
	testMinCostFlow2()
}
