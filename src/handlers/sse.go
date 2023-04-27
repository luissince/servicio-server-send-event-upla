package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Body struct {
	Id      string `json:"id,omitempty"`
	Codigo  string `json:"codigo"`
	Titulo  string `json:"titulo"`
	Mensaje string `json:"mensaje"`
}

type EventMessage struct {
	EventName string
	Data      any
}

type Client struct {
	Id          string
	SendMessage chan EventMessage
}

type HandlerEvent struct {
	mutex   sync.Mutex
	clients map[string]*Client
}

func NewHandlerEvent() *HandlerEvent {
	return &HandlerEvent{
		clients: make(map[string]*Client),
	}
}

func (h *HandlerEvent) addSubscription(id string, eventChan chan EventMessage) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.clients[id] = &Client{
		Id:          id,
		SendMessage: eventChan,
	}
	fmt.Println("Connected: ", id)
	fmt.Println("Clients: ", len(h.clients))
	fmt.Println("")
}

func (h *HandlerEvent) removeSubscription(id string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.clients, id)
	fmt.Println("Desconected: ", id)
	fmt.Println("Clients: ", len(h.clients))
	fmt.Println("")
}

func (h *HandlerEvent) sendNotificationById(id string, eventChan EventMessage) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	fmt.Println(h.clients)

	s, ok := h.clients[id]

	if ok {
		fmt.Println("enviando al usuario ", id, " ...")
		fmt.Println("")
		s.SendMessage <- eventChan
	}
}

func (h *HandlerEvent) sendNotificationAll(eventChan EventMessage) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	fmt.Println(h.clients)

	fmt.Println("enviando a todos...")
	fmt.Println("")

	for _, client := range h.clients {
		s, ok := h.clients[client.Id]
		if ok {
			s.SendMessage <- eventChan
		}
	}

}

func (h *HandlerEvent) HandlerSubcription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "El ID no existe", http.StatusBadRequest)
		return
	}

	eventChan := make(chan EventMessage)

	h.addSubscription(id, eventChan)

	defer h.removeSubscription(id)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Println("Streaming...")
	fmt.Println("")

	fmt.Fprintf(w, "data: Connected\n\n")
	flusher.Flush()

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				return
			}

			data, err := json.Marshal(event.Data)
			if err != nil {
				return
			}

			const format = "data:%s\n\n"
			fmt.Fprintf(w, format, string(data))
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

func (h *HandlerEvent) HandlerNotifyById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body := Body{}
	json.NewDecoder(r.Body).Decode(&body)

	event := EventMessage{
		EventName: "notification",
		Data:      body,
	}

	h.sendNotificationById(strings.ToUpper(body.Codigo), event)
	w.WriteHeader(http.StatusOK)
}

func (h *HandlerEvent) HandlerNotifyAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body := Body{}
	json.NewDecoder(r.Body).Decode(&body)

	id := uuid.New()

	body.Id = id.String()

	event := EventMessage{
		EventName: "notification",
		Data:      body,
	}

	h.sendNotificationAll(event)
	w.WriteHeader(http.StatusOK)
}
