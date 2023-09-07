package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// test if map from building song slice from file is not nil
func TestParseJSONToSlice(t *testing.T) {
	// Load the .env file using godotenv
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}

	songs, err := ParseJSONToSlice(os.Getenv("JSON_PATH"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if songs == nil {
		t.Error("Expected songs to be not nil")
	}
}

// test song path
func TestBuildSongPath(t *testing.T) {
	// Load the .env file using godotenv
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}

	songMap, err := ParseJSONToSlice(os.Getenv("JSON_PATH"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var songPath = BuildSongPath("Never Gonna Give You Up", songMap)
	try := os.Getenv("SONGS_PATH") + "/never_gonna_give_you_up.mp3"
	if songPath != try {
		t.Error("Expected song path to be", try, "got", songPath)
	}
}
