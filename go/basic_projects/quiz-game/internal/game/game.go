package game

import (
	"sync"
	"time"
)


type Game struct {
	Players []*Player
	Question []Question

	mu	sync.Mutex
}

func NewGame(players []*Player, questions []Question) *Game{
	return &Game{
		Players: players,
		Question: questions,
	}
}

func (g *Game) AskQuestion(q Question, duration time.Duration) {
	answers := make(chan Answer)

	for _, player := range g.Players{
		go ListenForAnswer(player, answers)
	}

	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case ans := <-answers:
			if ans.Text == q.Answer {
				g.mu.Lock()
				ans.Player.Score++
				g.mu.Unlock()
				return
			}
		case <-timer.C:
			return
		}
	}

}
