package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"servicio-sse-upla/src/handlers"

	"github.com/joho/godotenv"
)


func main() {
	time.LoadLocation("America/Lima")
	godotenv.Load()

	var go_port string = os.Getenv("GO_PORT")

	router := http.NewServeMux()
	handlers.InitRouters(router)
	fmt.Println("Ejecutando la aplicación:", go_port)
	http.ListenAndServe(go_port, router)
}
