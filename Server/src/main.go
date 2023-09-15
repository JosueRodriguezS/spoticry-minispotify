package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	http.Handle("/", r)
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