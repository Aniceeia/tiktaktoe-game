package services

import "math"

type AIService struct{}

func NewAI() *AIService {
	return &AIService{}
}

func (ai *AIService) FindBestMove(board [3][3]int) (int, int) {
	bestScore := math.MinInt32
	var move [2]int

	for i := range 3 {
		for j := range 3 {
			if board[i][j] == 0 {
				board[i][j] = 1
				score := minimax(board, 0, false)
				board[i][j] = 0

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
	if result != 0 {
		switch result {
		case 1:
			return 10 - depth
		case 2:
			return depth - 10
		default:
			return 0
		}
	}

	if isMaximizing {
		maxScore := math.MinInt32
		for i := range 3 {
			for j := range 3 {
				if board[i][j] == 0 {
					board[i][j] = 1
					score := minimax(board, depth+1, false)
					board[i][j] = 0
					maxScore = max(maxScore, score)
				}
			}
		}
		return maxScore
	}
	minScore := math.MaxInt32
	for i := range 3 {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				board[i][j] = 2
				score := minimax(board, depth+1, true)
				board[i][j] = 0
				minScore = min(minScore, score)
			}
		}
	}
	return minScore

}

func CheckTerminal(board [3][3]int) int {
	for i := range 3 {
		if board[i][0] != 0 && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return board[i][0]
		}
	}

	for j := range 3 {
		if board[0][j] != 0 && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return board[0][j]
		}
	}

	if board[0][0] != 0 && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0]
	}
	if board[0][2] != 0 && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2]
	}

	for i := range 3 {
		for j := range 3 {
			if board[i][j] == 0 {
				return 0
			}
		}
	}
	return -1
}
