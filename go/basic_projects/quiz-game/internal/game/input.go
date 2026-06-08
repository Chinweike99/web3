package game

import (
	"bufio"
	"os"
	"strings"
)


type Answer struct {
	Player *Player
	Text string
}


func ListenForAnswer(player *Player, ch chan <- Answer) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	text = strings.TrimSpace(text)

	ch <- Answer{
		Player: player,
		Text: text,
	}

}