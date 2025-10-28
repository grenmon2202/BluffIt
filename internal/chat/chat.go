package chat

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uint32
	Content   string
	SenderID  uint32
	Timestamp int64
}

type ChatStore struct {
	mu       sync.RWMutex
	Messages map[uint32]*Message
}

func CreateMessage(content string, senderID uint32) *Message {
	return &Message{
		ID:        uuid.New().ID(),
		Content:   content,
		SenderID:  senderID,
		Timestamp: time.Now().Unix(),
	}
}

func (cs *ChatStore) AddMessage(message *Message) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.Messages == nil {
		cs.Messages = make(map[uint32]*Message)
	}

	cs.Messages[message.ID] = message
}
