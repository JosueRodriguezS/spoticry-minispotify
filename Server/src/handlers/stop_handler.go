package handlers

import (
	"fmt"
	"net/http"
)

func StopHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "stop")
}
