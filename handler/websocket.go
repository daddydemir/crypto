package handler

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Broadcaster struct {
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

var hub = &Broadcaster{
	clients: make(map[*websocket.Conn]bool),
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error upgrading connection: ", err, "")
		return
	}
	defer ws.Close()

	hub.mutex.Lock()
	hub.clients[ws] = true
	hub.mutex.Unlock()

	slog.Info("new connection", "total", len(hub.clients))

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			hub.mutex.Lock()
			delete(hub.clients, ws)
			hub.mutex.Unlock()
			break
		}
	}
}

func ListenAndBroadcast(client *redis.Client) {
	ctx := context.Background()
	pubsub := client.Subscribe(ctx, "market:prices")
	ch := pubsub.Channel()

	for msg := range ch {
		hub.mutex.Lock()
		for clnt := range hub.clients {
			err := clnt.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				clnt.Close()
				delete(hub.clients, clnt)
			}
		}
		hub.mutex.Unlock()
	}
}
