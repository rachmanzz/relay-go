package event

import (
	"sync"
)

// Broadcaster manages active SSE connections and dispatches events.
type Broadcaster struct {
	mu          sync.RWMutex
	connections map[string]chan Message
}

type Message struct {
	Event string
	Data  string
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		connections: make(map[string]chan Message),
	}
}

// Subscribe adds a new connection to the broadcaster.
func (b *Broadcaster) Subscribe(id string) chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Message, 10)
	b.connections[id] = ch
	return ch
}

// Unsubscribe removes a connection from the broadcaster.
func (b *Broadcaster) Unsubscribe(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if ch, ok := b.connections[id]; ok {
		close(ch)
		delete(b.connections, id)
	}
}

// Broadcast sends a message to all connected clients.
func (b *Broadcaster) Broadcast(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.connections {
		select {
		case ch <- msg:
		default:
			// Buffer full, skip or handle accordingly
		}
	}
}

// SendTo sends a message to a specific client ID.
func (b *Broadcaster) SendTo(id string, msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if ch, ok := b.connections[id]; ok {
		select {
		case ch <- msg:
		default:
		}
	}
}
