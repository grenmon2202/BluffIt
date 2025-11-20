package broadcast

import (
	"sync"

	"github.com/grenmon2202/bluffit/ws/schema/response"
)

type BroadcastHub struct {
	mu      sync.RWMutex
	lobbies map[uint32]map[*BroadcastClient]bool
}

func CreateBroadcastHub() *BroadcastHub {
	return &BroadcastHub{
		lobbies: make(map[uint32]map[*BroadcastClient]bool),
	}
}

func (h *BroadcastHub) AddSubscriber(lobbyID uint32, client *BroadcastClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.lobbies[lobbyID]; !exists {
		h.lobbies[lobbyID] = make(map[*BroadcastClient]bool)
	}

	h.lobbies[lobbyID][client] = true
}

func (h *BroadcastHub) RemoveSubscriber(lobbyID uint32, client *BroadcastClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subscribers, exists := h.lobbies[lobbyID]; exists {
		delete(subscribers, client)
		if len(subscribers) == 0 {
			delete(h.lobbies, lobbyID)
		}
	}
}

func (h *BroadcastHub) BroadcastToLobby(lobbyID uint32, message response.BroadcastBody) {
	h.mu.RLock()
	subs, exists := h.lobbies[lobbyID]
	if !exists {
		h.mu.RUnlock()
		return
	}

	clients := make([]*BroadcastClient, 0, len(subs))
	for c := range subs {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		c.Send(message) // safe, single writer owns WriteJSON
	}
}
