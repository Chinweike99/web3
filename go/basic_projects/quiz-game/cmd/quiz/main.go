package main

import (
	"fmt"
	"quiz-game/internal/game"
	"time"
)


func main() {
	players := []*game.Player{
		{Name: "Player 1"},
		{Name: "Player 2"},
	}

	question := []game.Question{
		{Text: "What is 2 + 2?", Answer: "4"},
		{Text: "What is the capital of France?", Answer: "Paris"},
	}

	g := game.NewGame(players, question)
	for _, q := range g.Question {
		fmt.Println(q.Text)
		g.AskQuestion(q, 5*time.Second)
	}

	fmt.Println("\nFinal Scores")
	for _, p := range players {
		fmt.Printf("%s: %d\n", p.Name, p.Score)
	}


}
