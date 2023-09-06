package models

import (
	"fmt"
	"os"
	"sync"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	/*"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"*/)

type SongPlayer struct {
	//queue of songs
	Queue []Song

	//Current song
	current beep.StreamSeekCloser

	mutex sync.Mutex
}

// Function to open and decode the audio/mp3 file in our current playList
func OpenAndDecodeMp3File(songName string, songs map[int]Song) (beep.StreamSeekCloser, error) {
	//get song path
	filepath := BuildSongPath(songName, songs)
	//check if the file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		//retorna un error que indica que el path no existe
		return nil, fmt.Errorf("File does not exist at path: %s", filepath)
	}

	//open the audio file|mp3
	file, err := os.Open(filepath)
	if err != nil {
		//return an specific error for not being able to open the file
		return nil, fmt.Errorf("Error opening the file: %v", err)
	}

	//Decode the audio file|mp3 in a stream of audio using beep
	myStreamer, _, err := mp3.Decode(file) // Ignoramos la variable "format".
	if err != nil {
		file.Close()
		//return an specific error for not being able to decode the file
		return nil, fmt.Errorf("Error opening the file: %v", err)
	}

	//Return the streamer decoded
	return myStreamer, nil
}

/*
func Run(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, ready, err := oto.NewContext(d.SampleRate(), 2, 2)
	if err != nil {
		return err
	}
	<-ready

	p := c.NewPlayer(d)
	defer p.Close()
	p.Play()

	fmt.Printf("Length: %d[bytes]\n", d.Length())
	for {
		time.Sleep(time.Second)
		if !p.IsPlaying() {
			break
		}
	}

	return nil
}

// Function to play a song by its name
func PlaySong(name string, songs map[int]Song) {

	//get working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Get the song full path
	path := dir + GetSongPath(name, songs)

	print("the song relative path is: ")
	print(GetSongPath(name, songs))
	print("the song path is: ")
	print(path)

	// Play the song
	Run(path)
}
*/
