package gametree

import (
	"math"

	"github.com/JonasBernard/min-cost-max-flow/graph"
	"github.com/JonasBernard/min-cost-max-flow/graphlike"
)

type GameState[T any] interface {
	graphlike.RootedTreelike[GameState[T]]
	Evaluate() float64
	IsTerminal() bool
	ApplyMove(*graph.WeightedDirectedEdge[GameState[T]]) graph.WeigthedDirectedGraph[GameState[T]]
}

// AlphaBetaSearch performs the alpha-beta pruning algorithm on a game tree.
// It returns the best score and the corresponding move.
func AlphaBetaSearch[V graph.Node, T GameState[V]](g graph.WeigthedDirectedGraph[V], depth int, alpha, beta float64, maximizingPlayer bool) (float64, *graph.WeightedDirectedEdge[V]) {
	if depth == 0 || g.IsTerminal() {
		return g.Evaluate(), nil
	}

	var bestMove *graph.WeightedDirectedEdge[V]

	if maximizingPlayer {
		maxEval := math.Inf(-1)
		for _, edge := range g.Edges {
			childGraph := g.ApplyMove(edge)
			eval, _ := AlphaBetaSearch(childGraph, depth-1, alpha, beta, false)
			if eval > maxEval {
				maxEval = eval
				bestMove = edge
			}
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval, bestMove
	} else {
		minEval := math.Inf(1)
		for _, edge := range g.Edges {
			childGraph := g.ApplyMove(edge)
			eval, _ := AlphaBetaSearch(childGraph, depth-1, alpha, beta, true)
			if eval < minEval {
				minEval = eval
				bestMove = edge
			}
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval, bestMove
	}
}
