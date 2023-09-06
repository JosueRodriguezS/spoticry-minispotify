package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"
)

type Song struct {
	Name   string `json:"Name"`
	Artist string `json:"Artist"`
	Genre  string `json:"Genre"`
	Path   string `json:"FilePath"`
}

type SongList struct {
	songs map[int]Song
}

func ParseJSONToSongMap(filename string) (map[int]Song, error) {
	// Read the JSON file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Create a map to store songs
	songMap := make(map[int]Song)

	// Unmarshal the JSON data into a map
	if err := json.Unmarshal(data, &songMap); err != nil {
		return nil, err
	}

	return songMap, nil
}

func CreateSongList() *SongList {
	return &SongList{}
}

// Function to build a song path
func BuildSongPath(songName string, songs map[int]Song) string {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}
	relativePath := GetSongRelativePath(songName, songs)
	return os.Getenv("SONGS_PATH") + relativePath
}

func GetSongRelativePath(name string, songs map[int]Song) string {
	for _, song := range songs {
		if song.Name == name {
			return song.Path
		}
	}
	return ""
}

// Function to get a song by its name from the song map
func GetSongByName(name string, songs map[int]Song) (Song, bool) {
	for _, song := range songs {
		if song.Name == name {
			return song, true
		}
	}
	return Song{}, false
}
