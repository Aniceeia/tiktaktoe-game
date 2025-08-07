package service

import "math"

type AI struct{}

func NewAI() *AI {
	return &AI{}
}

func (ai *AI) FindBestMove(board [3][3]int) (int, int) {
	bestScore := math.MinInt32
	var move [2]int

	for i := range 3 {
		for j := range 3 {
			if board[i][j] == empty {
				board[i][j] = o
				score := minimax(board, 0, false)
				board[i][j] = empty

				if score > bestScore {
					bestScore = score
					move = [2]int{i, j}
				}
			}
		}
	}
	return move[0], move[1]
}

func minimax(board [3][3]int, depth int, isMaximizing bool) int {
	result := CheckTerminal(board)
	if result != empty {
		switch result {
		case o:
			return 10 - depth
		case x:
			return depth - 10
		default:
			return 0
		}
	}

	if isMaximizing {
		maxScore := math.MinInt32
		for i := range 3 {
			for j := range 3 {
				if board[i][j] == empty {
					board[i][j] = o
					score := minimax(board, depth+1, false)
					board[i][j] = empty
					maxScore = max(maxScore, score)
				}
			}
		}
		return maxScore
	}
	minScore := math.MaxInt32
	for i := range 3 {
		for j := 0; j < 3; j++ {
			if board[i][j] == empty {
				board[i][j] = x
				score := minimax(board, depth+1, true)
				board[i][j] = empty
				minScore = min(minScore, score)
			}
		}
	}
	return minScore

}

func CheckTerminal(board [3][3]int) int {
	for i := range 3 {
		if board[i][0] != empty && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return board[i][0]
		}
	}

	for j := range 3 {
		if board[0][j] != empty && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return board[0][j]
		}
	}

	if board[0][0] != empty && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0]
	}
	if board[0][2] != empty && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2]
	}

	for i := range 3 {
		for j := range 3 {
			if board[i][j] == empty {
				return empty
			}
		}
	}
	return -1
}
