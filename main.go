// package main

// import (
// 	// "fmt"
// 	"log"
// 	"net/http"
// 	"servicio-sse-upla/handlers"
// )

// func main() {
// 	r := http.NewServeMux()
// 	handlers.InitRouters(r)
// 	log.Println("0.0.0.0:9000")
// 	err := http.ListenAndServe("0.0.0.0:9000", r)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }

package main

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"servicio-sse-upla/src/handlers"

	"github.com/joho/godotenv"
)

type user struct {
	id   int
	name string
}

type notification struct {
	userID int
	event  string
	data   any
}

type subscription struct {
	user  *user
	event chan notification
}

type hub struct {
	subscriptions map[int]*subscription
	mu            sync.Mutex
}

func (h *hub) addSubscription(userID int, event chan notification) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscriptions == nil {
		h.subscriptions = make(map[int]*subscription)
	}
	h.subscriptions[userID] = &subscription{
		user:  &user{id: userID},
		event: event,
	}
}

func (h *hub) removeSubscription(userID int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscriptions == nil {
		return
	}
	delete(h.subscriptions, userID)
}

func (h *hub) sendNotification(n notification) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subscriptions == nil {
		return
	}
	if s, ok := h.subscriptions[n.userID]; ok {
		s.event <- n
	}
}

func main() {
	time.LoadLocation("America/Lima")
	godotenv.Load()

	var go_port string = os.Getenv("GO_PORT")

	router := http.NewServeMux()
	handlers.InitRouters(router)
	fmt.Println("Ejecutando la aplicación:", go_port)
	http.ListenAndServe(go_port, router)

	// time.LoadLocation("America/Lima")
	// godotenv.Load()

	// var go_port string = os.Getenv("GO_PORT")

	// h := &hub{}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
	// 	w.Header().Set("Content-Type", "text/event-stream")
	// 	w.Header().Set("Cache-Control", "no-cache")
	// 	w.Header().Set("Connection", "keep-alive")

	// 	if r.Method != "GET" {
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	userID := 1
	// 	eventChan := make(chan notification)

	// 	h.addSubscription(userID, eventChan)

	// 	defer h.removeSubscription(userID)

	// 	flusher, ok := w.(http.Flusher)
	// 	if !ok {
	// 		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)

	// 	fmt.Fprintf(w, "data: Connected\n\n")
	// 	flusher.Flush()
	// 	for {
	// 		select {
	// 		case n, ok := <-eventChan:
	// 			if !ok {
	// 				return
	// 			}

	// 			data, err := json.Marshal(n.data)
	// 			if err != nil {
	// 				return
	// 			}

	// 			const format = "data:%s\n\n"
	// 			fmt.Fprintf(w, format, string(data))
	// 			flusher.Flush()

	// 		case <-r.Context().Done():
	// 			return
	// 		}
	// 	}
	// })

	// http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != "POST" {
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	var body = map[string]any{}
	// 	json.NewDecoder(r.Body).Decode(&body)

	// 	userID := 1
	// 	event := "notification"
	// 	data := body
	// 	h.sendNotification(notification{userID: userID, event: event, data: data})
	// 	w.WriteHeader(http.StatusOK)
	// })

	// fmt.Println("Ejecutando la aplicación:", go_port)
	// http.ListenAndServe(go_port, nil)
}
