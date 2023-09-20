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

func ParseJSONToSlice(filename string) ([]Song, error) {
	var songs []Song
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return songs, err
	}
	fmt.Println("Successfully Opened", filename)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return songs, err
	}
	// we initialize our Users array
	err = json.Unmarshal(byteValue, &songs)
	if err != nil {
		fmt.Println(err)
		return songs, err
	}
	return songs, nil
}

func CreateSongList() *SongList {
	return &SongList{}
}

// Function to build a song path
func BuildSongPath(songName string, songs []Song) string {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}
	relativePath := GetSongRelativePath(songName, songs)
	return os.Getenv("SONGS_PATH") + relativePath
}

func GetSongRelativePath(name string, songs []Song) string {
	for _, song := range songs {
		if song.Name == name {
			return song.Path
		}
	}
	return ""
}

// Function to get a song by its name from the song map
func GetSongByName(name string, songs []Song) (Song, bool) {
	for _, song := range songs {
		if song.Name == name {
			return song, true
		}
	}
	return Song{}, false
}
