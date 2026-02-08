// cmd/dev-server/handlers/websocket.go
package handlers

import (
    "log"
    "net/http"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
)

type LiveReloadServer struct {
    clients    map[*websocket.Conn]bool
    clientsMu  sync.RWMutex
    upgrader   websocket.Upgrader
    pingPeriod time.Duration
}

func NewLiveReloadServer() *LiveReloadServer {
    return &LiveReloadServer{
        clients: make(map[*websocket.Conn]bool),
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool { return true },
        },
        pingPeriod: 30 * time.Second,
    }
}

func (s *LiveReloadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    conn, err := s.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    defer conn.Close()
    
    // Register client
    s.clientsMu.Lock()
    s.clients[conn] = true
    s.clientsMu.Unlock()
    
    // Start ping/pong to keep connection alive
    go s.keepAlive(conn)
    
    // Wait for client disconnect
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break
        }
    }
    
    // Unregister client
    s.clientsMu.Lock()
    delete(s.clients, conn)
    s.clientsMu.Unlock()
}

func (s *LiveReloadServer) BroadcastReload(reason string) {
    s.clientsMu.RLock()
    defer s.clientsMu.RUnlock()
    
    message := map[string]interface{}{
        "type":   "reload",
        "reason": reason,
        "time":   time.Now().Unix(),
    }
    
    for client := range s.clients {
        go func(c *websocket.Conn) {
            c.WriteJSON(message)
        }(client)
    }
}

