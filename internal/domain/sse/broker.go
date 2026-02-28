package sse

import (
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type SSEEvent struct {
	Event string      // "notification", "heartbeat" 등
	Data  interface{} // JSON 직렬화될 데이터
}

type Broker struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]map[chan SSEEvent]struct{}
	done    chan struct{}
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[uuid.UUID]map[chan SSEEvent]struct{}),
		done:    make(chan struct{}),
	}
}

// Done, 브로커가 종료될 때 닫히는 채널을 반환합니다.
func (b *Broker) Done() <-chan struct{} {
	return b.done
}

// Shutdown closes all client channels and signals done.
func (b *Broker) Shutdown() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for userID, set := range b.clients {
		for ch := range set {
			close(ch)
		}
		delete(b.clients, userID)
	}

	close(b.done)
}

func (b *Broker) Subscribe(userID uuid.UUID) chan SSEEvent {
	ch := make(chan SSEEvent, 16)

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.clients[userID] == nil {
		b.clients[userID] = make(map[chan SSEEvent]struct{})
	}
	b.clients[userID][ch] = struct{}{}
	slog.Debug("subscribe", "clients", b.clients[userID])

	return ch
}

func (b *Broker) Unsubscribe(userID uuid.UUID, ch chan SSEEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	slog.Debug("unsubscribe", "clients", b.clients[userID])

	if set, ok := b.clients[userID]; ok {
		delete(set, ch)
		close(ch)
		if len(set) == 0 {
			delete(b.clients, userID)
		}
	}
}

func (b *Broker) Publish(userID uuid.UUID, event SSEEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	set, ok := b.clients[userID]
	if !ok {
		return
	}

	for ch := range set {
		// non-blocking send
		select {
		case ch <- event:
		default: // 버퍼가 가득 차있을 때 블락되는 걸 방지
		}
	}
}
