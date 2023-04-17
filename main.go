// package main

// import (
// 	"io"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-contrib/sse"
// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()

// 	// Middleware para CORS
// 	config := cors.DefaultConfig()
// 	config.AllowOrigins = []string{"*"}
// 	r.Use(cors.New(config))

// 	r.GET("/stream", func(c *gin.Context) {
// 		c.Stream(func(w io.Writer) bool {
// 			// Enviar un solo evento y luego cerrar la conexión
// 			data := sse.Event{
// 				Data: []byte("Hello world!"),
// 			}
// 			c.SSEvent("message", data)
// 			return false // La conexión se cierra automáticamente después de enviar el evento
// 		})
// 	})

// 	r.Run(":9000")
// }

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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

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

	h := &hub{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userID := 1
		eventChan := make(chan notification)
		h.addSubscription(userID, eventChan)
		defer h.removeSubscription(userID)
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "data: Connected\n\n")
		flusher.Flush()
		for {
			select {
			case n, ok := <-eventChan:
				if !ok {
					return
				}

				data, err := json.Marshal(n.data)
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
	})

	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body = map[string]any{}
		json.NewDecoder(r.Body).Decode(&body)

		userID := 1
		event := "notification"
		data := body
		h.sendNotification(notification{userID: userID, event: event, data: data})
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Ejecutando la aplicación:", go_port)
	http.ListenAndServe(go_port, nil)
}

// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// type User struct {
// 	Name     string
// 	Messages []string
// 	sync.Mutex
// }

// var users = map[string]*User{}
// var mutex = sync.Mutex{}

// func main() {
// 	http.HandleFunc("/sse", handleSSE)
// 	http.HandleFunc("/notify", handleNotify)
// 	http.ListenAndServe("0.0.0.0:9000", nil)
// }

// func handleSSE(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")

// 	userID := r.URL.Query().Get("user")
// 	if userID == "" {
// 		userID = fmt.Sprintf("user-%d", time.Now().Unix())
// 	}

// 	user := getUser(userID)
// 	defer user.Unlock()

// 	user.Lock()

// 	fmt.Fprintf(w, "event: welcome\n")
// 	fmt.Fprintf(w, "data: Bienvenido %s!\n\n", user.Name)

// 	for _, message := range user.Messages {
// 		fmt.Fprintf(w, "event: message\n")
// 		fmt.Fprintf(w, "data: %s\n\n", message)
// 	}

// 	user.Messages = []string{}

// 	notify := w.(http.CloseNotifier).CloseNotify()
// 	go func() {
// 		<-notify
// 		mutex.Lock()
// 		delete(users, userID)
// 		mutex.Unlock()
// 	}()

// 	for {
// 		// time.Sleep(10 * time.Second)
// 		fmt.Fprintf(w, "event: ping\n")
// 		fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
// 	}
// }

// func handleNotify(w http.ResponseWriter, r *http.Request) {
// 	// userID := r.FormValue("user")
// 	// message := r.FormValue("message")

// 	userID := "76423388"
// 	message := "hola perro"

// 	if userID == "" || message == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	user := getUser(userID)
// 	// defer user.Unlock()

// 	user.Lock()
// 	defer user.Unlock()
// 	user.Messages = append(user.Messages, message)
// 	// user.Unlock()
// }

// func getUser(userID string) *User {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	if user, ok := users[userID]; ok {
// 		return user
// 	}

// 	user := &User{Name: userID}
// 	users[userID] = user
// 	return user
// }

// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// type User struct {
// 	sync.Mutex
// 	ID       string
// 	Messages []string
// }

// var users = make(map[string]*User)
// var mutex = &sync.Mutex{}

// func main() {
// 	http.HandleFunc("/", handleIndex)
// 	http.HandleFunc("/sse", handleSSE)
// 	http.HandleFunc("/notify", handleNotify)

// 	fmt.Println("Servidor escuchando en http://localhost:9000")
// 	http.ListenAndServe("0.0.0.0:9000", nil)
// }

// func handleIndex(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Bienvenido al servidor de eventos con Go")
// }

// func handleSSE(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")

// 	userID := r.URL.Query().Get("user")
// 	if userID == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println("userID: ", userID)

// 	user := getUser(userID)

// 	// Set headers for server-sent events
// 	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	// w.Header().Set("Access-Control-Allow-Origin", "*")
// 	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
// 	// w.Header().Set("Content-Type", "text/event-stream")
// 	// w.Header().Set("Cache-Control", "no-cache")
// 	// w.Header().Set("Connection", "keep-alive")

// 	// Continuously send messages to the client
// 	for {
// 		// Lock the user to avoid data races
// 		user.Lock()

// 		// If there are messages, send them to the client
// 		if len(user.Messages) > 0 {
// 			message := user.Messages[0]
// 			user.Messages = user.Messages[1:]
// 			fmt.Fprintf(w, "data: %s\n\n", message)
// 		}

// 		// Unlock the user to allow other goroutines to access it
// 		user.Unlock()

// 		// Sleep for a second to avoid spamming the client
// 		time.Sleep(1 * time.Second)
// 	}
// }

// func handleNotify(w http.ResponseWriter, r *http.Request) {
// 	userID := r.FormValue("user")
// 	message := r.FormValue("message")

// 	if userID == "" || message == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println("userID: ", userID)
// 	fmt.Println("message: ", message)

// 	user := getUser(userID)

// 	// Lock the user to avoid data races
// 	user.Lock()

// 	// Add the message to the user's message queue

// 	fmt.Println("Send ok")
// 	user.Messages = append(user.Messages, message)

// 	// Unlock the user to allow other goroutines to access it
// 	user.Unlock()

// 	// Return a success status
// 	w.WriteHeader(http.StatusOK)
// }

// func getUser(id string) *User {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	if user, ok := users[id]; ok {
// 		return user
// 	}

// 	user := &User{ID: id}
// 	users[id] = user
// 	return user
// }
