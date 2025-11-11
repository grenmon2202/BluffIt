package requests

import "encoding/json"

type Envelope struct {
	Type      string          `json:"Type"`
	RequestID string          `json:"RequestID"`
	Data      json.RawMessage `json:"Data"`
}

type CreateLobbyRequest struct {
	LobbyHost string `json:"lobbyHost"`
}

type FetchLobbyRequest struct {
	LobbyID uint32 `json:"lobbyId"`
}

type JoinLobbyRequest struct {
	LobbyCode  string `json:"lobbyCode"`
	PlayerName string `json:"playerName"`
}

type SendMessageRequest struct {
	LobbyID  uint32 `json:"lobbyId"`
	SenderID uint32 `json:"senderId"`
	Content  string `json:"content"`
}
