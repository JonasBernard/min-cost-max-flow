package minimax

import (
	"math"
)

const (
	MaxPlayer = 1
	MinPlayer = -1
)

type GameState[T any] interface {
	IsTerminal() bool
	GeneratePossibleMoves(bool) []T
	ValueOfTerminal() int
}

func MinimaxMinimize[T GameState[T]](state GameState[T]) (bestMove GameState[T]) {
	bestMove, _ = MinimaxRecurse[T](state, false)
	return
}

func MinimaxMaximize[T GameState[T]](state GameState[T]) (bestMove GameState[T]) {
	bestMove, _ = MinimaxRecurse[T](state, true)
	return
}

func MinimaxRecurse[T GameState[T]](state GameState[T], isMaximizing bool) (GameState[T], int) {
	if state.IsTerminal() {
		return nil, state.ValueOfTerminal()
	}

	bestVal := math.MaxInt
	if isMaximizing {
		bestVal = math.MinInt
	}
	bestMove := state
	for _, move := range state.GeneratePossibleMoves(!isMaximizing) {
		_, val:= MinimaxRecurse[T](move, !isMaximizing)
		if (val > bestVal && isMaximizing) || (val < bestVal && !isMaximizing) {
			bestVal = val
			bestMove = move
		}
	}
	return bestMove, bestVal
}

