package events

import (
	"fmt"
	"net/http"
	"sync"
)

type EventMessage struct {
	EventName string
	Data      any
}

type HandlerEvent struct {
	m       sync.Mutex
	clients map[string]*client
}

func NewHandlerEvent() *HandlerEvent {
	return &HandlerEvent{
		clients: make(map[string]*client),
	}
}

func (h *HandlerEvent) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id := r.URL.Query().Get("id")

	if id == "" {
		fmt.Println("Error id")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := newClient(id)
	h.register(c)
	fmt.Println("Connected: ", id)
	c.onLine(r.Context(), w, flusher)
	fmt.Println("Desconected: ", id)
	h.removeClient(id)
}

func (h *HandlerEvent) register(c *client) {
	h.m.Lock()
	defer h.m.Unlock()
	h.clients[c.Id] = c
}

func (h *HandlerEvent) removeClient(id string) {
	h.m.Lock()
	defer delete(h.clients, id)
	h.m.Unlock()
}

func (h *HandlerEvent) Broadcast(m EventMessage) {
	h.m.Lock()
	defer h.m.Unlock()
	for _, c := range h.clients {
		c.sendMessage <- m
	}
}
