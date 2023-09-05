package models

import (
	"encoding/json"
	"io/ioutil"
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
