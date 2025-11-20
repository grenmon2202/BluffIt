package lobby

import (
	"math/rand"
	"strings"
	"sync"

	"github.com/grenmon2202/bluffit/internal/player"
)

type LobbyStore struct {
	mu      sync.RWMutex
	lobbies map[uint32]*Lobby
}

func NewLobbyStore() *LobbyStore {
	return &LobbyStore{
		lobbies: make(map[uint32]*Lobby),
	}
}

func (lobbyStore *LobbyStore) AddLobby(lobbyHost *player.Player) Lobby {
	lobby := CreateLobby(lobbyHost, lobbyStore.GenerateUniqueLobbyCode())

	lobbyStore.mu.Lock()
	lobbyStore.lobbies[lobby.ID] = lobby
	lobbyStore.mu.Unlock()

	return *lobby
}

func (lobbyStore *LobbyStore) RemoveLobbyByID(lobbyID uint32) {
	lobbyStore.mu.Lock()
	defer lobbyStore.mu.Unlock()

	lobby := lobbyStore.lobbies[lobbyID]
	if lobby != nil {
		lobby.LobbyDeleted = true
	}

	delete(lobbyStore.lobbies, lobbyID)
}

func (lobbyStore *LobbyStore) GetAllLobbies() map[uint32]Lobby {
	lobbyStore.mu.RLock()
	defer lobbyStore.mu.RUnlock()

	lobbies := make(map[uint32]Lobby)
	for id, lobby := range lobbyStore.lobbies {
		lobbies[id] = *lobby
	}
	return lobbies
}

func (lobbyStore *LobbyStore) GetLobbyByID(lobbyID uint32) (*Lobby, bool) {
	lobbyStore.mu.RLock()
	defer lobbyStore.mu.RUnlock()

	lobby, exists := lobbyStore.lobbies[lobbyID]
	if !exists {
		return nil, false
	}

	return lobby, exists
}

func (lobbyStore *LobbyStore) GetLobbyByCode(code string) (*Lobby, bool) {
	lobbyStore.mu.RLock()
	defer lobbyStore.mu.RUnlock()

	for _, lobby := range lobbyStore.lobbies {
		if lobby.UniqueCode == code {
			return lobby, true
		}
	}

	return nil, false
}

func (lobbyStore *LobbyStore) GenerateUniqueLobbyCode() string {
	const codeLength = 8
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for {
		var code strings.Builder
		for i := 0; i < codeLength; i++ {
			randomIndex := rand.Intn(len(charset))
			code.WriteByte(charset[randomIndex])
		}

		generatedCode := code.String()
		if _, exists := lobbyStore.GetLobbyByCode(generatedCode); !exists {
			return generatedCode
		}
	}
}

func (lobbyStore *LobbyStore) JoinLobbyByCode(code string, name string) (bool, string, player.Player, Lobby) {
	lobbyStore.mu.RLock()
	defer lobbyStore.mu.RUnlock()

	var foundLobby *Lobby

	for _, lobby := range lobbyStore.lobbies {
		if lobby.UniqueCode == code {
			foundLobby = lobby
			break
		}
	}

	if foundLobby == nil {
		return false, "Lobby not found", player.Player{}, Lobby{}
	}

	if foundLobby.Players.CheckPlayerExistsByName(name) {
		return false, "Player name already taken in this lobby", player.Player{}, Lobby{}
	}

	newPlayer := player.CreatePlayer(name, false)
	foundLobby.Players.AddPlayer(newPlayer)

	return true, "Successfully joined lobby", *newPlayer, *foundLobby
}
