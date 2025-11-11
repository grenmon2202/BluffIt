package lobby

import (
	"math/rand"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/grenmon2202/bluffit/internal/chat"
	"github.com/grenmon2202/bluffit/internal/player"
)

// Lobby holds connected players and metadata for a game session.
type Lobby struct {
	ID          uint32
	Players     *player.PlayerStore
	Chats       *chat.ChatStore
	UniqueCode  string
	LobbyHostID uint32
}

type LobbyStore struct {
	mu      sync.RWMutex
	lobbies map[uint32]*Lobby
}

func NewLobbyStore() *LobbyStore {
	return &LobbyStore{
		lobbies: make(map[uint32]*Lobby),
	}
}

func (lobby *Lobby) AddChatMessage(content string, senderID uint32) (bool, string) {
	if !lobby.Players.CheckPlayerExistsByID(senderID) {
		return false, "Sender does not exist in the lobby"
	}

	message := chat.CreateMessage(content, senderID)

	lobby.Chats.AddMessage(message)
	return true, "Message added successfully"
}

func (lobbyStore *LobbyStore) AddLobby(lobbyHost *player.Player) Lobby {
	lobby := &Lobby{
		ID:          uuid.New().ID(),
		Players:     &player.PlayerStore{},
		Chats:       &chat.ChatStore{},
		UniqueCode:  lobbyStore.GenerateUniqueLobbyCode(),
		LobbyHostID: lobbyHost.ID,
	}

	lobby.Players.AddPlayer(lobbyHost)

	lobbyStore.mu.Lock()
	lobbyStore.lobbies[lobby.ID] = lobby
	lobbyStore.mu.Unlock()

	return *lobby
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

func (lobbyStore *LobbyStore) JoinLobbyByCode(code string, name string) (bool, string, player.Player) {
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
		return false, "Lobby not found", player.Player{}
	}

	if foundLobby.Players.CheckPlayerExistsByName(name) {
		return false, "Player name already taken in this lobby", player.Player{}
	}

	newPlayer := player.CreatePlayer(name, false)
	foundLobby.Players.AddPlayer(newPlayer)

	return true, "Successfully joined lobby", *newPlayer
}
