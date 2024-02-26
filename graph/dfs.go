package graph

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
