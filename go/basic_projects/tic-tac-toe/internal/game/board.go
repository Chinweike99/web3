package game

import "fmt"

type Board struct {
	Cells [9]rune
}

func NewBoard() *Board {
	b := &Board{}
	for i := range b.Cells {
		b.Cells[i] = ' '
	}
	return  b
}


func (b *Board) PlaceMove(pos int, mark rune) bool {
	if pos < 0 || pos >= 9 {
		fmt.Println("Invalid position")
		return  false
	}
	if b.Cells[pos] != ' '{
		return  false
	}
	b.Cells[pos] = mark
	return true
}


func (b *Board) IsFull() bool {
	for _, c := range b.Cells {
		if c == ' ' {
			return  false
		}
	}
	return  true
}

func (b *Board) Winner()rune {
	winPatterns := [8][3]int {
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	for _, p := range winPatterns {
		if b.Cells[p[0]] != ' ' &&
			b.Cells[p[0]] == b.Cells[p[1]] &&
			b.Cells[p[1]] == b.Cells[p[2]] {
				return  b.Cells[p[0]]
		}
	}
	return  ' '
}

