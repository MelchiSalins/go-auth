package handler

import (
	"fmt"
	"net/http"
)

// RootPath handle request on / path
func RootPath(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Root Path")
}
