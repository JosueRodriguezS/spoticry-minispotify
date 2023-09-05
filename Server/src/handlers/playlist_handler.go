package handlers

import (
	"fmt"
	"net/http"
)

func PlaylistHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "playlist")
}
