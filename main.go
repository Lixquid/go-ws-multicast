package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
	verbose bool
}

func (h *Hub) add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn] = struct{}{}
}

func (h *Hub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, conn)
}

func (h *Hub) broadcast(sender *websocket.Conn, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.verbose {
		log.Printf("broadcast: %s", message)
	}

	for client := range h.clients {
		if client == sender {
			continue
		}

		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("write error: %v", err)
		}
	}
}

func handler(hub *Hub) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade error: %v", err)
			return
		}
		defer conn.Close()

		hub.add(conn)
		defer hub.remove(conn)

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if messageType == websocket.TextMessage {
				hub.broadcast(conn, message)
			}
		}
	})

	return mux
}

func main() {
	host := flag.String("host", "localhost", "host to listen on")
	port := flag.Int("port", 9994, "port to listen on")
	verbose := flag.Bool("verbose", false, "print every message sent through the relay")
	flag.Parse()

	hub := &Hub{
		clients: make(map[*websocket.Conn]struct{}),
		verbose: *verbose,
	}

	addr := *host + ":" + strconv.Itoa(*port)
	log.Printf("ws-multicast listening on ws://%s", addr)

	if err := http.ListenAndServe(addr, handler(hub)); err != nil {
		log.Fatal(err)
	}
}

