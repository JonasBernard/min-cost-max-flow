package graph

import "errors"

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
			return nil, errors.New("no path from root to vertex")
		}
		path.Vertices = append(path.Vertices, parent)
		path.Edges = append(path.Edges, p.getEdge(parent, head))
		head = parent
	}
}
