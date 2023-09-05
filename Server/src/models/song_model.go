package models

import (
	"encoding/json"
	"io/ioutil"
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

func GetSongPath(name string, songs map[int]Song) string {
	for _, song := range songs {
		if song.Name == name {
			return song.Path
		}
	}
	return ""
}
