package main

import (
	"io/ioutil"
	"net/http"
)

func main() {

	// Create a file server which serves files out of the "./static" directory.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/getAudioBuffer", func(w http.ResponseWriter, r *http.Request) {
		// Read the audio buffer from a file or generate it dynamically
		audioBuffer, err := ioutil.ReadFile("C:/Users/josue/OneDrive/Documents/TEC/Semestre_II_2023/Leguajes/proyecto/testVLC/src/ric_king.mp3")
		if err != nil {
			http.Error(w, "Unable to read audio buffer", http.StatusInternalServerError)
			return
		}

		// Set the appropriate content type for the response
		w.Header().Set("Content-Type", "audio/mpeg")

		// Write the audio buffer to the response
		w.Write(audioBuffer)
	})

	println("Hello, World!")
	http.ListenAndServe(":8080", nil)
}
