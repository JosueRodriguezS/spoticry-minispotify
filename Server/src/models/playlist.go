package models

type Playlist struct {
	
	Id int;
	Name string;
	Songs []Song;

}
// Constructor
func NewPlaylist(id int, name string, songs []Song) *Playlist {
	return &Playlist{
		Id: id,
		Name: name,
		Songs: songs,
	}
}

// Agregar canciones al playlis
func (p *Playlist) AddSong(s Song) {
	p.Songs = append(p.Songs, s)
}

// Eliminar canciones del playlist
func (p *Playlist) DeleteSong(s Song) {
	for i, song := range p.Songs {
		if song.Name == s.Name {
			p.Songs = append(p.Songs[:i], p.Songs[i+1:]...)
		}
	}
}