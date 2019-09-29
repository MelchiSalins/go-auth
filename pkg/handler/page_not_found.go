package handler

import "net/http"

// CustomPageNotFound when a 404 is generated
func CustomPageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Custom 404 Page"))
}
