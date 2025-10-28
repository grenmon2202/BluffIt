package game

import "github.com/google/uuid"

// Game represents the state and logic of an ongoing game session.
type Game struct {
	ID uint32
}

func newGame() *Game {
	return &Game{
		ID: uuid.New().ID(),
	}
}
