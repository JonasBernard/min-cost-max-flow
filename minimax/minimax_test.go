package minimax_test

import (
	"fmt"
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/minimax"
)

type TicTacToe struct {
	gameboard [][]int
}

func (t TicTacToe) String() string {
	first := fmt.Sprintf("%v%v%v", t.gameboard[0][0], t.gameboard[0][1], t.gameboard[0][2])
	second := fmt.Sprintf("%v%v%v", t.gameboard[1][0], t.gameboard[1][1], t.gameboard[1][2])
	third := fmt.Sprintf("%v%v%v", t.gameboard[2][0], t.gameboard[2][1], t.gameboard[2][2])
	term := "N"
	if t.IsTerminal() {
		term = "T"
	}
	return fmt.Sprintf("%v<%v|%v|%v>%v", term, first, second, third, t.ValueOfTerminal())
}

func (t TicTacToe) IsTerminal() bool {
	if t.ValueOfTerminal() != 0 {
		return true
	}

	return t.NumberOfEmptyFields() == 0
}

func (t TicTacToe) NumberOfEmptyFields() int {
	count := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t.gameboard[i][j] == 0 {
				count++
			}
		}
	}
	return count
}

func (t TicTacToe) GeneratePossibleMoves(forMinimizer bool) []TicTacToe {
	possibleMoves := make([]TicTacToe, 0)
	if t.IsTerminal() {
		return possibleMoves
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t.gameboard[i][j] == 0 {
				newGameboard := make([][]int, 3)
				for i := range newGameboard {
					newGameboard[i] = make([]int, 3)
					copy(newGameboard[i], t.gameboard[i])
				}
				if forMinimizer {
					newGameboard[i][j] = 2
				} else {
					newGameboard[i][j] = 1
				}
				possibleMoves = append(possibleMoves, TicTacToe{gameboard: newGameboard})
			}
		}
	}
	return possibleMoves
}

func (t TicTacToe) winnerValue(i, j int) int {
	X := 1
	O := 2
	if t.gameboard[i][j] == X {
		return t.NumberOfEmptyFields() + 1
	} else if t.gameboard[i][j] == O {
		return -t.NumberOfEmptyFields() - 1
	}
	return 0
}

func (t TicTacToe) ValueOfTerminal() (value int) {
	value = 0

	for i := 0; i < 3; i++ {
		// check rows
		if t.gameboard[i][0] == t.gameboard[i][1] && t.gameboard[i][1] == t.gameboard[i][2] {
			value = t.winnerValue(i, 0)
		}
		// check columns
		if t.gameboard[0][i] == t.gameboard[1][i] && t.gameboard[1][i] == t.gameboard[2][i] {
			value = t.winnerValue(0, i)
		}
	}

	// check diagonals
	if t.gameboard[0][0] == t.gameboard[1][1] && t.gameboard[1][1] == t.gameboard[2][2] {
		value = t.winnerValue(0, 0)
	}
	if t.gameboard[0][2] == t.gameboard[1][1] && t.gameboard[1][1] == t.gameboard[2][0] {
		value = t.winnerValue(0, 2)
	}

	return
}

func TestTicTacToe(t *testing.T) {
	X := 1
	O := 2

	gameboard := [][]int{
		[]int{X, 0, 0},
		[]int{0, 0, O},
		[]int{0, O, O},
	}

	game := TicTacToe{gameboard: gameboard}

	fmt.Println(minimax.MinimaxMaximize[TicTacToe](game))
}
