package market

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Rooms map[string]map[*websocket.Conn]struct{}
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) Add(symbol string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	symbol = strings.ToUpper(symbol)

	if h.Rooms[symbol] == nil {
		h.Rooms[symbol] = make(map[*websocket.Conn]struct{})
	}

	h.Rooms[symbol][conn] = struct{}{}
}

func (h *Hub) Remove(symbol string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	symbol = strings.ToUpper(symbol)

	if room, ok := h.Rooms[symbol]; ok {
		delete(room, conn)

		if len(room) == 0 {
			delete(h.Rooms, symbol)
		}
	}

	conn.Close()
}

func (h *Hub) Broadcast(symbol string, payload any) {
	h.mu.RLock()
	room, ok := h.Rooms[strings.ToUpper(symbol)]
	if !ok {
		h.mu.RUnlock()
		return
	}

	conns := make([]*websocket.Conn, 0, len(room))
	for c := range room {
		conns = append(conns, c)
	}
	h.mu.RUnlock()

	b, _ := json.Marshal(payload)

	for _, c := range conns {
		if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
			go h.Remove(symbol, c)
		}
	}
}