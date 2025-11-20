package player

import (
	"github.com/google/uuid"
)

// Player represents a single participant in the game.
type Player struct {
	ID     uint32
	Name   string
	IsHost bool
}

func CreatePlayer(name string, isHost bool) *Player {
	return &Player{
		ID:     uuid.New().ID(),
		Name:   name,
		IsHost: isHost,
	}
}
