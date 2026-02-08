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

// keepAlive sends periodic ping messages to keep the WebSocket connection alive
func (s *LiveReloadServer) keepAlive(conn *websocket.Conn) {
    ticker := time.NewTicker(s.pingPeriod)
    defer ticker.Stop()
    
    // Configure connection for pong handling
    conn.SetPongHandler(func(appData string) error {
        // Reset read deadline on pong
        conn.SetReadDeadline(time.Now().Add(s.pingPeriod * 2))
        return nil
    })
    
    for {
        select {
        case <-ticker.C:
            // Set write deadline for ping
            conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                // If ping fails, remove client and exit goroutine
                log.Printf("Failed to send ping: %v", err)
                s.removeClient(conn)
                return
            }
            
            // Set read deadline expecting pong
            conn.SetReadDeadline(time.Now().Add(s.pingPeriod * 2))
        }
    }
}

// removeClient safely removes a client from the connections map
func (s *LiveReloadServer) removeClient(conn *websocket.Conn) {
    s.clientsMu.Lock()
    defer s.clientsMu.Unlock()
    
    if _, exists := s.clients[conn]; exists {
        delete(s.clients, conn)
        conn.Close()
    }
}

func (s *LiveReloadServer) BroadcastReload(reason string) {
    s.clientsMu.RLock()
    defer s.clientsMu.RUnlock()
    
    if len(s.clients) == 0 {
        return
    }
    
    message := map[string]interface{}{
        "type":   "reload",
        "reason": reason,
        "time":   time.Now().Unix(),
    }
    
    for client := range s.clients {
        go func(c *websocket.Conn) {
            // Set write deadline for broadcast
            c.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.WriteJSON(message); err != nil {
                log.Printf("Failed to send reload message: %v", err)
                s.removeClient(c)
            }
        }(client)
    }
}

// BroadcastMessage sends a custom JSON message to all connected clients
func (s *LiveReloadServer) BroadcastMessage(messageType string, data interface{}) {
    s.clientsMu.RLock()
    defer s.clientsMu.RUnlock()
    
    if len(s.clients) == 0 {
        return
    }
    
    message := map[string]interface{}{
        "type": messageType,
        "data": data,
        "time": time.Now().Unix(),
    }
    
    for client := range s.clients {
        go func(c *websocket.Conn) {
            c.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.WriteJSON(message); err != nil {
                log.Printf("Failed to send message: %v", err)
                s.removeClient(c)
            }
        }(client)
    }
}

// GetClientCount returns the number of currently connected clients
func (s *LiveReloadServer) GetClientCount() int {
    s.clientsMu.RLock()
    defer s.clientsMu.RUnlock()
    return len(s.clients)
}

// CloseAll closes all WebSocket connections
func (s *LiveReloadServer) CloseAll() {
    s.clientsMu.Lock()
    defer s.clientsMu.Unlock()
    
    for client := range s.clients {
        client.Close()
    }
    s.clients = make(map[*websocket.Conn]bool)
}
