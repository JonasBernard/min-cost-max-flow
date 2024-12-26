package gametree

import (
    "math"
)

type GameState[T any] interface {
    RootedTreelike[GameState[T]]
	Evaluate() float64
}

// AlphaBetaSearch performs the alpha-beta pruning algorithm on a game tree.
// It returns the best score and the corresponding move.
func AlphaBetaSearch[T GameState](g WeigthedDirectedGraph[T], depth int, alpha, beta float64, maximizingPlayer bool) (float64, *WeightedDirectedEdge[T]) {
    if depth == 0 || g.IsTerminal() {
        return g.Evaluate(), nil
    }

    var bestMove *WeightedDirectedEdge[T]

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
