package lobby

import (
	"github.com/grenmon2202/bluffit/internal/player"
)

// Lobby holds connected players and metadata for a game session.
type Lobby struct {
	ID           uint32
	Players      *player.PlayerStore
	UniqueCode   string
	LobbyHostID  uint32
	GameStarted  bool
	LobbyDeleted bool
}

func CreateLobby(lobbyHost *player.Player, uniqueCode string) *Lobby {
	lobby := &Lobby{
		ID:          lobbyHost.ID, // Temporary ID assignment; should be replaced with a proper unique ID generator
		Players:     player.CreatePlayerStore(),
		UniqueCode:  uniqueCode,
		LobbyHostID: lobbyHost.ID,
	}

	lobby.Players.AddPlayer(lobbyHost)

	return lobby
}
