package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHub_ClientLifecycleAndBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWS(w, r)
	}))
	t.Cleanup(func() {
		server.Close()
	})

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// 1. Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket server: %v", err)
	}
	t.Cleanup(func() {
		ws.Close()
	})

	// Allow registration loop to run
	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	clientCount := len(hub.clients)
	hub.mu.RUnlock()
	if clientCount != 1 {
		t.Fatalf("Expected 1 registered WebSocket client, got %d", clientCount)
	}

	// 2. Broadcast JSON event
	type TestEvent struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}
	event := TestEvent{Type: "test_event", Data: "hello_swarm"}
	hub.BroadcastJSON(event)

	// Read message from client WS
	var received TestEvent
	ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err := ws.ReadJSON(&received); err != nil {
		t.Fatalf("Failed to read broadcast JSON from WebSocket client: %v", err)
	}

	if received.Type != "test_event" || received.Data != "hello_swarm" {
		t.Errorf("Unexpected event received: %+v", received)
	}

	// 3. Close client connection & verify unregistration
	ws.Close()
	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	clientCountAfter := len(hub.clients)
	hub.mu.RUnlock()
	if clientCountAfter != 0 {
		t.Errorf("Expected 0 clients after disconnect, got %d", clientCountAfter)
	}
}
