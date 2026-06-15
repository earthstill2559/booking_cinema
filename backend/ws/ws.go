package ws

import (
	"net/http"
	"sync"

	"cinema/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clientsMu sync.RWMutex
	clients   = make(map[*websocket.Conn]bool)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func Broadcast(showtimeID, seatID, status, lockedBy string) {
	msg := models.WSMessage{
		Type:       "seat_update",
		SeatID:     seatID,
		ShowtimeID: showtimeID,
		Status:     status,
		LockedBy:   lockedBy,
	}

	clientsMu.RLock()
	conns := make([]*websocket.Conn, 0, len(clients))
	for conn := range clients {
		conns = append(conns, conn)
	}
	clientsMu.RUnlock()

	var stale []*websocket.Conn
	for _, conn := range conns {
		if err := conn.WriteJSON(msg); err != nil {
			stale = append(stale, conn)
		}
	}

	if len(stale) > 0 {
		clientsMu.Lock()
		for _, conn := range stale {
			delete(clients, conn)
			conn.Close()
		}
		clientsMu.Unlock()
	}
}
