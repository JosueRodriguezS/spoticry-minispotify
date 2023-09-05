package handlers

import (
	"fmt"
	"net/http"
)

func PlaylistHandler(w http.ResponseWriter, r *http.Request) {
	playlist := "My Number One\neye of the tiger\nnever gonna give you up\nnever gonna give you up\nric king\nsandstorm\nthrough the fire and flames\nwe are number one"
	fmt.Fprintln(w, playlist)
}
