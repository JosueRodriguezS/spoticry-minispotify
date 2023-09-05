// print play

package handlers

import (
	"fmt"
	"net/http"
)

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "play")
}
