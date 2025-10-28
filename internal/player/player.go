package player

import (
	"sync"

	"github.com/google/uuid"
)

// Player represents a single participant in the game.
type Player struct {
	ID     uint32
	Name   string
	IsHost bool
}

// PlayerStore maintains a global registry of all players
type PlayerStore struct {
	mu      sync.RWMutex
	Players map[uint32]*Player
}

func CreatePlayer(name string, isHost bool) *Player {
	return &Player{
		ID:     uuid.New().ID(),
		Name:   name,
		IsHost: isHost,
	}
}

func (ps *PlayerStore) AddPlayer(player *Player) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.Players == nil {
		ps.Players = make(map[uint32]*Player)
	}

	ps.Players[player.ID] = player
}

func (ps *PlayerStore) CheckPlayerExistsByID(id uint32) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, exists := ps.Players[id]
	return exists
}

func (ps *PlayerStore) CheckPlayerExistsByName(name string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, player := range ps.Players {
		if player.Name == name {
			return true
		}
	}
	return false
}
