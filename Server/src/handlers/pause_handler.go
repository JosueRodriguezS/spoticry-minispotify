package handlers

import (
	"fmt"
	"net/http"
)

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pause")
}
