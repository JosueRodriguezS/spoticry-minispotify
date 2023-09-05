package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Song struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Genre  string `json:"genre"`
	Path   string `json:"path"`
}

type SongList struct {
	songs map[string]Song
}

func populateMap(jsonpath string) map[string]Song {
	// Open our jsonFile
	jsonFile, err := os.Open(jsonpath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + jsonpath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)
	// we initialize our Users array
	var songs []Song
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &songs)
	// json structure = {id:{name: "name", artist: "artist", genre: "genre", path: "relativepath = ./songs/"}}
	songsMap := make(map[string]Song)
	for _, song := range songs {
		songsMap[song.Name] = song
	}
	return songsMap

}
func CreateSongList() *SongList {
	return &SongList{}
}
