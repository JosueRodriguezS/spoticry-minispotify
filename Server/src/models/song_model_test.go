package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// test if map from building song map from file is not nil
func TestParseJSONToSongMap(t *testing.T) {
	// Load the .env file using godotenv
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error al cargar el archivo .env:", err)
	}

	var songs, ready = ParseJSONToSongMap(os.Getenv("JSON_PATH"))
	fmt.Println("THIS IS MY SONGS PATH")
	fmt.Println(os.Getenv("JSON_PATH"))
	if songs == nil {
		//print ready
		fmt.Println(ready)
		t.Error("Expected map of songs, got nil")
	}
}
