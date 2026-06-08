package game


type Game struct {
	Board *Board
	Players [2]Player
	Turn int
}

func NewGame(p1, p2 Player) *Game {
	return  &Game{
		Board: NewBoard(),
		Players: [2]Player{p1, p2},
		Turn: 0,
	}
}

func (g *Game) CurrentPlayer() Player {
	return  g.Players[g.Turn%2]
}

func (g *Game) MakeMove(pos int) bool {
	player := g.CurrentPlayer()
	if g.Board.PlaceMove(pos, player.Mark){
		g.Turn++
		return true
	}
	return false
}

func (g *Game) IsOver() (bool, rune) {
	if winner := g.Board.Winner(); winner != ' ' {
		return  true, winner
	}
	if g.Board.IsFull() {
		return true, 'D'
	}
	return false, ' '
}




