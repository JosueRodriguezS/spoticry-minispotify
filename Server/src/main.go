package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Song struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	FilePath string `json:"filePath"`
	Artist   string `json:"artist"`
}

var songs []Song

func main() {

	// Load the .env file using godotenv
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}

	loadSongs()

	r := mux.NewRouter()
	r.HandleFunc("/songs", getSongList).Methods("GET")
	r.HandleFunc("/songs", addSong).Methods("POST")
	r.HandleFunc("/songs/{Name}", getSong).Methods("GET")

	// UpdateAddSong
	r.HandleFunc("/addSong", AddSongJson).Methods("POST")
	r.HandleFunc("/delete/{Name}", deleteSong).Methods("POST")
	r.HandleFunc("/getBuffer/{Name}", getBuffer).Methods("GET")

	//Funciones para buscar canciones
	r.HandleFunc("/songs/firstLetter/{Letter}", getSongsByFirstLetter).Methods("GET")
	r.HandleFunc("/songs/wordcount/{Count}", searchSongsByWordCount).Methods("GET")
	r.HandleFunc("/songs/fileSizeRange/{minSize}", searchSongsByFileSizeRange).Methods("GET")

	// Crear un middleware CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Cambia esto según la URL de tu cliente web
		AllowedMethods: []string{"GET", "POST"},           // Los métodos permitidos
	})
	handler := corsMiddleware.Handler(r)

	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}

func loadSongs() {
	file, err := os.Open(os.Getenv("SONGS_PATH") + "/songs.json")
	if err != nil {
		log.Fatal("Error opening songs.json:", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading songs.json:", err)
		return
	}

	err = json.Unmarshal(data, &songs)
	if err != nil {
		log.Fatal("Error unmarshaling songs.json:", err)
		return
	}
}

func getSongList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func getSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Log the incoming request
	fmt.Println("Request for song:", params["Name"])

	for _, song := range songs {
		if song.Name == params["Name"] {
			json.NewEncoder(w).Encode(song)
			return
		}
	}

	// Return a proper error response if the song is not found
	http.Error(w, "Song not found", http.StatusNotFound)
}

func addSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON request body into a Song struct
	var newSong Song
	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if a song with the same name already exists
	for _, existingSong := range songs {
		if existingSong.Name == newSong.Name {
			http.Error(w, "Song with the same name already exists", http.StatusConflict)
			return
		}
	}

	// Add the new song to the songs slice
	songs = append(songs, newSong)

	// Return the newly added song as the response
	json.NewEncoder(w).Encode(newSong)
}

func getBuffer(w http.ResponseWriter, r *http.Request) {
	// Obtener el nombre del archivo MP3 de los parámetros
	params := mux.Vars(r)
	songName := params["Name"]

	// Construir la ruta completa al archivo MP3 utilizando la variable FilePath en la estructura Song
	var filePath string
	for _, song := range songs {
		if song.Name == songName {
			filePath = os.Getenv("SONGS_PATH") + "/" + song.FilePath
			break
		}
	}
	if filePath == "" {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}
	// Leer el archivo MP3 en bytes
	mp3Bytes, err := ioutil.ReadFile(filePath)
	fmt.Println(filePath)
	if err != nil {
		http.Error(w, "Error reading MP3 file", http.StatusInternalServerError)
		return
	}
	// Configurar las cabeceras de la respuesta HTTP
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Length", string(len(mp3Bytes)))

	// Escribir los bytes del archivo MP3 en la respuesta HTTP
	_, err = w.Write(mp3Bytes)
	if err != nil {
		http.Error(w, "Error writing MP3 data to response", http.StatusInternalServerError)
		return
	}

}

// UpdateAddSong
func AddSongJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newSong Song

	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if a song with the same name already exists
	for _, existingSong := range songs {
		if existingSong.Name == newSong.Name {
			http.Error(w, "Song with the same name already exists", http.StatusConflict)
			return
		}
	}

	// Add the new song to the songs slice
	songs = append(songs, newSong)

	// Save the updated list of songs to songs.json
	if err := saveSongsToFile(); err != nil {
		http.Error(w, "Error saving songs to file", http.StatusInternalServerError)
		return
	}

	// Return the newly added song as the response
	json.NewEncoder(w).Encode(newSong)
}

func saveSongsToFile() error {
	// Encode the list of songs to JSON
	encodedData, err := json.Marshal(songs)
	if err != nil {
		return err
	}

	// Write the JSON data to the songs.json file
	if err := ioutil.WriteFile(os.Getenv("SONGS_PATH")+"/songs.json", encodedData, 0644); err != nil {
		return err
	}

	return nil
}

func deleteSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	nameToDelete := params["Name"]

	// Busca la canción por su nombre y la elimina de la lista
	for i, song := range songs {
		if song.Name == nameToDelete {
			// Elimina la canción de la lista
			songs = append(songs[:i], songs[i+1:]...)

			// Guarda la lista actualizada de canciones en songs.json
			if err := saveSongsToFile(); err != nil {
				http.Error(w, "Error saving songs to file", http.StatusInternalServerError)
				return
			}

			// Devuelve una respuesta exitosa
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Si no se encontró la canción, devuelve un error 404
	http.Error(w, "Song not found", http.StatusNotFound)
}

// Funcion para buscar canciones por la primera letra
func getSongsByFirstLetter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	firstLetter := params["Letter"]

	var matchingSongs []Song

	for _, song := range songs {
		if strings.HasPrefix(strings.ToLower(song.Name), strings.ToLower(firstLetter)) {
			matchingSongs = append(matchingSongs, song)
		}
	}

	if len(matchingSongs) == 0 {
		http.Error(w, "No songs found with the specified first letter", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(matchingSongs)
}

// Funcion para buscar canciones por numero de palabras
func searchSongsByWordCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	countStr := params["Count"]

	count, err := strconv.Atoi(countStr)
	if err != nil {
		http.Error(w, "Invalid count parameter", http.StatusBadRequest)
		return
	}

	matchingSongs := []Song{}
	for _, song := range songs {
		// Divide el nombre de la canción en palabras
		words := strings.Fields(song.Name)

		// Comprueba si el número de palabras coincide con el valor especificado
		if len(words) == count {
			matchingSongs = append(matchingSongs, song)
		}
	}

	if len(matchingSongs) == 0 {
		http.Error(w, "No songs found with the specified word count", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(matchingSongs)
}

// Funcion para buscar canciones por tamaño de archivo en megabytes
func searchSongsByFileSizeRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	minSizeMBStr := params["minSize"]
	

	minSizeMB, err := strconv.Atoi(minSizeMBStr)
	if err != nil {
		http.Error(w, "Invalid minSize parameter", http.StatusBadRequest)
		return
	}

	maxSizeMB := minSizeMB + 2

	matchingSongs := []Song{}
	for _, song := range songs {
		// Obtén el nombre del archivo MP3 de la canción
		songFilePath := os.Getenv("SONGS_PATH") + "/" + song.FilePath

		// Lee el tamaño del archivo MP3
		fileInfo, err := os.Stat(songFilePath)
		if err != nil {
			// Maneja el error si no se puede acceder al archivo
			continue
		}

		// Convierte el tamaño del archivo a MB
		fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024) // 1 MB = 1024 * 1024 bytes

		// Compara el tamaño del archivo con el rango especificado en MB
		if fileSizeMB >= float64(minSizeMB) && fileSizeMB <= float64(maxSizeMB) {
			matchingSongs = append(matchingSongs, song)
		}
	}

	if len(matchingSongs) == 0 {
		http.Error(w, "No songs found within the specified file size range", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(matchingSongs)
}
