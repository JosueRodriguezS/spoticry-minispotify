package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Song struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	FilePath string `json:"filePath"`
	Artist   string `json:"artist"`
}

type SearchRequest struct {
	Action   string `json:"action"`
	Genre    string `json:"genre"`
	Artist   string `json:"artist"`
	SongName string `json:"songName"`
}

var songs []Song

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	// Load the .env file using godotenv
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}

	loadSongs()

	//Create a new router
	r := mux.NewRouter()

	//Define the routes/endpoints and their handlers
	r.HandleFunc("/songs", getSongList).Methods("GET")
	r.HandleFunc("/songs", addSong).Methods("POST")
	r.HandleFunc("/songs/{Name}", getSong).Methods("GET")
	r.HandleFunc("/search", searchSongs).Methods("GET")

	//Start the HTTP server on port 8080
	go func() {
		http.Handle("/", r)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// Start your WebSocket server in a goroutine
	go func() {
		http.HandleFunc("/ws", handleWebSocket)
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatal("WebSocket server error:", err)
		}
	}()

	// Keep the main function running
	select {}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Cliente conectado")

	for {
		// Read message from client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			break
		}

		// Decode the received message
		var message map[string]interface{}
		if err := json.Unmarshal(p, &message); err != nil {
			fmt.Println("Error decoding WebSocket message:", err)
			break
		}

		// Check the action specified in the message
		action, ok := message["action"].(string)
		if !ok {
			fmt.Println("Invalid action in WebSocket message")
			break
		}

		switch action {
		case "GetSong":
			// Extract the song name from the message
			songName, ok := message["songName"].(string)
			if !ok {
				fmt.Println("Invalid songName in WebSocket message")
				break
			}

			// Find the song by name and send it as a response
			var foundSong Song
			for _, song := range songs {
				if song.Name == songName {
					foundSong = song
					break
				}
			}

			// Send the found song as a response
			response, err := json.Marshal(foundSong)
			if err != nil {
				fmt.Println("Error encoding song response:", err)
				break
			}

			// Send the response back to the client
			if err := conn.WriteMessage(messageType, response); err != nil {
				fmt.Println("Error sending song response:", err)
				break
			}
		case "SearchSong":
			searchRequest := SearchRequest{}
			if err := json.Unmarshal(p, &searchRequest); err != nil {
				fmt.Println("Error decoding WebSocket message:", err)
				break
			}

			// Perform the song search and get the search results
			searchResults := performSongSearch(searchRequest)

			// Send the search results as a response
			response, err := json.Marshal(searchResults)
			if err != nil {
				fmt.Println("Error encoding search response:", err)
				break
			}

			// Send the response back to the client
			if err := conn.WriteMessage(messageType, response); err != nil {
				fmt.Println("Error sending search response:", err)
				break
			}
		case "AddSong":

		default:
			fmt.Println("Unsupported action:", action)
		}
	}
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

func searchSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Parse the HTTP request parameters
	queryParams := r.URL.Query()
	genre := queryParams.Get("genre")
	artist := queryParams.Get("artist")
	name := queryParams.Get("name")

	// Filter songs by genre and artist
	var result []Song
	for _, song := range songs {
		if (genre == "" || song.Genre == genre) && (artist == "" || song.Artist == artist) && (name == "" || song.Name == name) {
			result = append(result, song)
		}
	}

	// Retrun the filtered songs as the response
	json.NewEncoder(w).Encode(result)
}

func performSongSearch(request SearchRequest) []Song {
	// Initialize an empty result slice
	var searchResults []Song

	// Iterate through your songs and filter based on the search criteria
	for _, song := range songs {
		// Check if the genre matches (if provided)
		if request.Genre != "" && song.Genre != request.Genre {
			continue
		}

		// Check if the artist matches (if provided)
		if request.Artist != "" && song.Artist != request.Artist {
			continue
		}

		// Check if the song name contains the search term (if provided)
		if request.SongName != "" && !strings.Contains(song.Name, request.SongName) {
			continue
		}

		// If all criteria match or no criteria provided, add the song to results
		searchResults = append(searchResults, song)
	}

	return searchResults
}

// Function to update the songs.json file by adding songs
func updateAddSong(newSongs []Song) error {
	// Open the songs.json file for reading and writing
	filePath := os.Getenv("SONGS_PATH") + "/songs.json"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening songs.json: %v", err)
	}
	defer file.Close()

	// Read the existing content
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading songs.json: %v", err)
	}

	// Unmarshal the existing content
	var songs []Song
	if err := json.Unmarshal(data, &songs); err != nil {
		return fmt.Errorf("error unmarshaling songs.json: %v", err)
	}

	// Append the new songs
	songs = append(songs, newSongs...)

	// Marshal the updated content
	updatedData, err := json.Marshal(songs)
	if err != nil {
		return fmt.Errorf("error marshaling songs.json: %v", err)
	}

	// Write the updated content back to the file
	if err := ioutil.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing songs.json: %v", err)
	}

	return nil
}

// Function to update the songs.json file by deleting songs
func updateDeleteSong(songsToDelete []Song) error {
	// Open the songs.json file for reading and writing
	filePath := os.Getenv("SONGS_PATH") + "/songs.json"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening songs.json: %v", err)
	}
	defer file.Close()

	// Read the existing content
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading songs.json: %v", err)
	}

	// Unmarshal the existing content
	var songs []Song
	if err := json.Unmarshal(data, &songs); err != nil {
		return fmt.Errorf("error unmarshaling songs.json: %v", err)
	}

	// Create a map of song names to be deleted for efficient lookup
	songsToDeleteMap := make(map[string]struct{})
	for _, song := range songsToDelete {
		songsToDeleteMap[song.Name] = struct{}{}
	}

	// Filter out the songs to be deleted
	var updatedSongs []Song
	for _, song := range songs {
		_, exists := songsToDeleteMap[song.Name]
		if !exists {
			updatedSongs = append(updatedSongs, song)
		}
	}

	// Marshal the updated content
	updatedData, err := json.Marshal(updatedSongs)
	if err != nil {
		return fmt.Errorf("error marshaling songs.json: %v", err)
	}

	// Write the updated content back to the file
	if err := ioutil.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing songs.json: %v", err)
	}

	return nil
}
