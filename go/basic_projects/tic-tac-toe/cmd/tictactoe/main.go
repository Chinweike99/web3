package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tic-tac-toe/internal/game"
)



func printBoard(b *game.Board) {
	fmt.Println()
	for i := 0; i < 9; i += 3 {
		fmt.Printf(" %c | %c | %c \n", b.Cells[i], b.Cells[i+1], b.Cells[i+2])
		if i < 6 {
			fmt.Println("---+---+---")
		}
	}
	fmt.Println()
}


func main() {
	reader := bufio.NewReader(os.Stdin)

	p1 := game.Player{Name: "Player 1", Mark: 'X'}
	p2 := game.Player{Name: "Player 2", Mark: 'O'}

	g := game.NewGame(p1, p2)
	for {
		printBoard(g.Board)
		player := g.CurrentPlayer()

		fmt.Printf("%s (%c), choose positon (1-9): ", player.Name, player.Mark)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		pos, err := strconv.Atoi(input)
		if err != nil || pos < 1 || pos > 9 {
			fmt.Println("Invalid input. Use numbers 1-9.")
			continue
		}

		if !g.MakeMove(pos - 1) {
			fmt.Println("That spot is already taken")
			continue
		}

		done, winner := g.IsOver()
		if done {
			printBoard(g.Board)
			if winner == 'D' {
				fmt.Println("It's a draw.")
			}else {
				fmt.Printf("Player %c wins!\n", winner)
			}
			break
		}
	}

}




