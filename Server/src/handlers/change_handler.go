package handlers

import (
	"fmt"
	"net/http"
)

func NextHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "next")
}

func PreviousHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "previous")
}
