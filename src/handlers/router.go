package handlers

import (
	"encoding/json"
	"net/http"
)

type Reponse struct {
	Message string `json:"message"`
}

func InitRouters(r *http.ServeMux) {
	handlerEvents := NewHandlerEvent()

	r.HandleFunc("/notify", handlerEvents.HandlerSubcription)
	r.HandleFunc("/senduser", handlerEvents.HandlerNotifyById)
	r.HandleFunc("/sendall", handlerEvents.HandlerNotifyAll)
	r.HandleFunc("/", Handler)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	reponse := Reponse{
		Message: "Server Send Event",
	}
	jsonResponse, _ := json.Marshal(reponse)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
