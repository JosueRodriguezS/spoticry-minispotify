package main

// API Rest to manage the server
import (
	"fmt"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"test.z/handlers"
	//"github.com/faiface/beep"
	//"github.com/faiface/beep/mp3"
	//"github.com/faiface/beep/speaker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/play", handlers.PlayHandler).Methods("GET")
	r.HandleFunc("/pause", handlers.PauseHandler).Methods("GET")
	r.HandleFunc("/stop", handlers.StopHandler).Methods("GET")
	r.HandleFunc("/next", handlers.NextHandler).Methods("GET")
	r.HandleFunc("/previous", handlers.PreviousHandler).Methods("GET")
	r.HandleFunc("/playlist", handlers.PlaylistHandler).Methods("GET")
	http.Handle("/", r)

	fmt.Println("Servidor API en ejecución en el puerto 8080...")
	fmt.Println(os.Getenv("SongsPath"))
	http.ListenAndServe(":8080", nil)
}

/* data := []byte("Hola, esto es un ejemplo de codificación Base64 en Go")
encoded := base64.StdEncoding.EncodeToString(data)
fmt.Println(encoded)

decoded, err := base64.StdEncoding.DecodeString(encoded)
if err != nil {
	fmt.Println("Error:", err)
	return
}
fmt.Println(string(decoded)) */
