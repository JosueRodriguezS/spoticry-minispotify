// print play

package handlers

import (
	"fmt"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gorilla/mux"
	"test.z/models"

	"net/http"
)

func PlayHandler(w http.ResponseWriter, r *http.Request, songMap map[int]models.Song) {
	// Use Gorilla mux to get the song name from the URL
	vars := mux.Vars(r)
	songName := vars["songName"]
	fmt.Fprintln(w, "play", songName)

	go func() {
		Mystreamer, err := models.OpenAndDecodeMp3File("My Number One", songMap)
		if err != nil {
			// Manejar errores aquí, si es necesario.
			fmt.Println("Error:", err)
			return
		}
		defer Mystreamer.Close()

		// Configurar el sistema de sonido utilizando speaker.Init con una configuración de formato predeterminada.
		sampleRate := beep.SampleRate(44100)
		speaker.Init(sampleRate, sampleRate.N(time.Second/10))

		// Iniciar la reproducción utilizando speaker.Play.
		speaker.Play(Mystreamer)

		// Mantener la goroutine en ejecución para que la canción se reproduzca.
		select {}
	}()
}
