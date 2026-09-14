package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestMulticast(t *testing.T) {
	hub := &Hub{
		clients: make(map[*websocket.Conn]struct{}),
	}

	server := httptest.NewServer(handler(hub))
	defer server.Close()

	// Convert the HTTP test server URL to a WebSocket URL.
	wsURL := "ws" + server.URL[len("http"):]

	client1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("client1 failed to connect: %v", err)
	}
	defer client1.Close()

	client2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("client2 failed to connect: %v", err)
	}
	defer client2.Close()

	// Give the server goroutines a moment to register both clients.
	time.Sleep(10 * time.Millisecond)

	message := "hello, clients"

	if err := client1.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	// Client 2 should receive the message.
	client2.SetReadDeadline(time.Now().Add(time.Second))

	messageType, received, err := client2.ReadMessage()
	if err != nil {
		t.Fatalf("client2 failed to receive message: %v", err)
	}

	if messageType != websocket.TextMessage {
		t.Fatalf("expected text message, got message type %d", messageType)
	}

	if string(received) != message {
		t.Fatalf("expected %q, got %q", message, received)
	}

	// Client 1 should NOT receive its own message.
	client1.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

	_, _, err = client1.ReadMessage()
	if err == nil {
		t.Fatal("client1 unexpectedly received its own message")
	}
}

