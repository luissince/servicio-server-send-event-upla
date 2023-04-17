package handlers

import (
	"encoding/json"
	"net/http"
	"servicio-sse-upla/handlers/events"
)

type Reponse struct {
	Message string `json:"message"`
}

func InitRouters(r *http.ServeMux) {
	handlerEvents := events.NewHandlerEvent()

	r.HandleFunc("/notify", handlerEvents.Handler)
	r.HandleFunc("/test1", HandlerTest1(handlerEvents))
	r.HandleFunc("/test2", HandlerTest2(handlerEvents))
	// r.Handle("/", http.FileServer(http.Dir("./public")))
	r.HandleFunc("/", Handler)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	reponse := Reponse{
		Message: "Server Send Event",
	}

	jsonResponse, err := json.Marshal(reponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func HandlerTest1(notifier *events.HandlerEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]any{}
		json.NewDecoder(r.Body).Decode(&data)
		notifier.Broadcast(events.EventMessage{
			EventName: "saludar",
			Data:      data,
		})
	}
}

func HandlerTest2(notifier *events.HandlerEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data = map[string]any{}
		json.NewDecoder(r.Body).Decode(&data)
		notifier.Broadcast(events.EventMessage{
			EventName: "saltar",
			Data:      data,
		})
	}
}
