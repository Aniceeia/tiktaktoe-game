package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidMove(t *testing.T) {
	tests := []struct {
		name     string
		oldBoard [3][3]int
		newBoard [3][3]int
		expected bool
	}{
		{
			"Valid move - place X",
			[3][3]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			[3][3]int{
				{1, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			true,
		},
		{
			"Invalid move - multiple changes",
			[3][3]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			[3][3]int{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			false,
		},
		{
			"Invalid move - overwrite existing",
			[3][3]int{
				{1, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			[3][3]int{
				{2, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidMove(tt.oldBoard, tt.newBoard)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckTerminal(t *testing.T) {
	tests := []struct {
		name     string
		board    [3][3]int
		expected int
	}{
		{
			"Player X wins horizontal",
			[3][3]int{
				{1, 1, 1},
				{0, 0, 0},
				{0, 0, 0},
			},
			1,
		},
		{
			"Player O wins vertical",
			[3][3]int{
				{2, 0, 0},
				{2, 0, 0},
				{2, 0, 0},
			},
			2,
		},
		{
			"Player X wins diagonal",
			[3][3]int{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			1,
		},
		{
			"Draw",
			[3][3]int{
				{1, 2, 1},
				{1, 2, 2},
				{2, 1, 1},
			},
			-1,
		},
		{
			"Game continues",
			[3][3]int{
				{1, 0, 0},
				{0, 2, 0},
				{0, 0, 1},
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckTerminal(tt.board)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindBestMove(t *testing.T) {
	ai := NewAI()

	tests := []struct {
		name     string
		board    [3][3]int
		expected [2]int
	}{
		{
			"Block opponent win",
			[3][3]int{
				{1, 0, 1},
				{0, 2, 0},
				{0, 0, 0},
			},
			[2]int{0, 1},
		},
		{
			"Win immediately",
			[3][3]int{
				{2, 0, 2},
				{0, 1, 0},
				{0, 0, 1},
			},
			[2]int{0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col := ai.FindBestMove(tt.board)
			assert.Equal(t, tt.expected[0], row)
			assert.Equal(t, tt.expected[1], col)
		})
	}
}
