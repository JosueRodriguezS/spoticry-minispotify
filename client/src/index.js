import React, { useState, useEffect, useRef } from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlay, faStop, faForward, faBackward, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';

function App() {
  const [searchType, setSearchType] = useState('firstLetter');
  const [searchValue, setSearchValue] = useState('');
  const [selectedSong, setSelectedSong] = useState('');
  const [canciones, setCanciones] = useState([]);
  const [playlists, setPlaylists] = useState([]);
  const [currentPlaylist, setCurrentPlaylist] = useState({ nombre: '', canciones: [] });
  const audioRef = useRef(null);

  useEffect(() => {
    getCanciones();
  }, []);

  function createPlaylist(nombre) {
    const newPlaylist = { nombre, canciones: [] };
    setPlaylists([...playlists, newPlaylist]);
  }

  function addSongToPlaylist(playlistIndex, song) {
    const updatedPlaylist = { ...playlists[playlistIndex] };
    if (!updatedPlaylist.canciones.includes(song)) {
      updatedPlaylist.canciones = [...updatedPlaylist.canciones, song];
      const updatedPlaylists = [...playlists];
      updatedPlaylists[playlistIndex] = updatedPlaylist;
      setPlaylists(updatedPlaylists);
    }
  }

  function removeSongFromPlaylist(playlistIndex, song) {
    const updatedPlaylist = { ...playlists[playlistIndex] };
    updatedPlaylist.canciones = updatedPlaylist.canciones.filter(c => c !== song);
    const updatedPlaylists = [...playlists];
    updatedPlaylists[playlistIndex] = updatedPlaylist;
    setPlaylists(updatedPlaylists);
  }

  function removePlaylist(playlistIndex) {
    const updatedPlaylists = playlists.filter((_, index) => index !== playlistIndex);
    setPlaylists(updatedPlaylists);
  }
  

  function getCanciones() {
    fetch('http://localhost:8080/songs')
      .then(response => response.json())
      .then(data => {
        setCanciones(data);
      })
      .catch(error => {
        console.error('Error al cargar canciones:', error);
      });
  }

  function search() {
    // Realiza la búsqueda según el tipo seleccionado (FirstLetter o WordCount)
    fetch(`http://localhost:8080/songs/${searchType}/${encodeURIComponent(searchValue)}`)
      .then(response => response.json())
      .then(data => {
        setCanciones(data);
      })
      .catch(error => {
        console.error('Error al realizar la búsqueda:', error);
      });
  }

  function play(songName) {
    // Hacer la solicitud GET para obtener el buffer del servidor
    fetch(`http://localhost:8080/getBuffer/${encodeURIComponent(songName)}`)
      .then(response => response.arrayBuffer())
      .then(data => {
        // Crear una URL para el buffer de audio
        const audioBlob = new Blob([data], { type: 'audio/mpeg' });
        const audioUrl = URL.createObjectURL(audioBlob);

        // Detener la reproducción actual antes de cargar una nueva canción
        if (!audioRef.current.paused) {
          audioRef.current.pause();
          audioRef.current.currentTime = 0;
        }

        // Reproducir la nueva canción
        audioRef.current.src = audioUrl;
        audioRef.current.play();
      })
      .catch(error => {
        console.error('Error al obtener el buffer de audio:', error);
      });
  }

  function playPlaylist(playlistIndex) {
    const playlist = playlists[playlistIndex];
    let currentSongIndex = 0;
  
    // Función para reproducir la siguiente canción en la lista
    function playNextSong() {
      if (currentSongIndex < playlist.canciones.length) {
        play(playlist.canciones[currentSongIndex]);
        currentSongIndex++;
      } else {
        // La lista de reproducción ha terminado
        currentSongIndex = 0;
      }
    }
  
    // Reproducir la primera canción al iniciar la lista de reproducción
    playNextSong();
  
    // Agregar un evento para detectar cuando una canción ha terminado de reproducirse
    audioRef.current.addEventListener("ended", () => {
      // Reproducir la siguiente canción
      playNextSong();
    });
  }
  
  
  

  function stop() {
    // Detener la reproducción actual
    audioRef.current.pause();
    audioRef.current.currentTime = 0;
  }

  function forward() {
    // Avanzar 10 segundos (ajusta según tus necesidades)
    audioRef.current.currentTime += 10;
  }

  function backward() {
    // Retroceder 10 segundos (ajusta según tus necesidades)
    audioRef.current.currentTime -= 10;
  }

  function getPlaceholderText(searchType) {
    switch (searchType) {
      case 'firstLetter':
        return 'la primera letra';
      case 'wordcount':
        return 'el número de palabras';
      case 'fileSizeRange':
        return 'el tamaño mínimo en MB';
      default:
        return '';
    }
  }
  return (
    <div className="container">
      <div>
        <h1>Buscar Canciones</h1>
        {/* Selector para elegir el tipo de búsqueda */}
        <select value={searchType} onChange={e => setSearchType(e.target.value)}>
          <option value="firstLetter">Primera Letra</option>
          <option value="wordcount">Numero de Palabras</option>
          <option value="fileSizeRange">Peso del archivo</option>
        </select>
        {/* Campo de entrada de texto para el valor de búsqueda */}
        <input
          type="text"
          placeholder={`Ingrese ${getPlaceholderText(searchType)}`}
          value={searchValue}
          onChange={e => setSearchValue(e.target.value)}
        />
        {/* Botón para realizar la búsqueda */}
        <button onClick={search}>Buscar</button>
        <button onClick={getCanciones}>Todas las canciones</button>
      </div>
      <h1>Listado de Canciones</h1>
      <div id="canciones">
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {canciones.map(cancion => (
            <li key={cancion.name} style={{ marginBottom: '10px' }}>
              <button onClick={() => play(cancion.name)}>
                <FontAwesomeIcon icon={faPlay} />
              </button>
              <span style={{ marginLeft: '10px' }}>{cancion.name}</span>
            </li>
          ))}
        </ul>
      </div>
      {/* Elemento <audio> para reproducir el audio */}
      <audio ref={audioRef} controls style={{ marginTop: '10px' }}>
        Tu navegador no soporta la reproducción de audio.
      </audio>
      <div className="player">
        <button onClick={backward}>
          <FontAwesomeIcon icon={faBackward} /> 
        </button>
        <button onClick={stop}>
          <FontAwesomeIcon icon={faStop} /> 
        </button>
        <button onClick={forward}>
          <FontAwesomeIcon icon={faForward} /> 
        </button>
      </div>
      <div>
        <h2>Listas de Reproducción</h2>
        <div>
        <input
            type="text"
            placeholder="Nombre de la lista"
            value={currentPlaylist.nombre}
            onChange={e => setCurrentPlaylist({ ...currentPlaylist, nombre: e.target.value })}
          />
          <button onClick={() => createPlaylist(currentPlaylist.nombre)}>
            <FontAwesomeIcon icon={faPlus} /> Crear Lista
          </button>
        </div>
        <ul>
          {playlists.map((playlist, index) => (
            <li key={index}>
              {playlist.nombre}
              <button onClick={() => removePlaylist(index)}>
                <FontAwesomeIcon icon={faTrash} /> Eliminar
              </button>
              <button onClick={() => playPlaylist(index)}>
                <FontAwesomeIcon icon={faPlay} /> Play Playlist
              </button>
              <ul>
                {playlist.canciones.map(song => (
                  <li key={song}>
                    {song}
                    <button onClick={() => removeSongFromPlaylist(index, song)}>
                      <FontAwesomeIcon icon={faTrash} /> Eliminar Canción
                    </button>
                    <button onClick={() => play(song)}>
                      <FontAwesomeIcon icon={faPlay} /> Play
                    </button>
                  </li>
                ))}
              </ul>
              <select onChange={e => setSelectedSong(e.target.value)}>
                <option value="">Seleccionar Canción</option>
                {canciones.map(song => (
                  <option key={song.name} value={song.name}>
                    {song.name}
                  </option>
                ))}
              </select>
              <button onClick={() => addSongToPlaylist(index, selectedSong)}>
                Agregar a Lista
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

ReactDOM.createRoot(document.getElementById('root')).render(<App />);
