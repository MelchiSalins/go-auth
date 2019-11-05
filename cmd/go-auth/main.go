package main

import (
	"log"
	"net/http"

	"github.com/MelchiSalins/go-auth/pkg/app"
	"github.com/MelchiSalins/go-auth/pkg/handler"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting server on port: ", app.Port)
	log.Fatal(http.ListenAndServe(app.Port, Handler()))
}

// Handler returns Gorilla Mux Handler for HTTP Server
func Handler() *mux.Router {
	auth, err := handler.NewAuthenticator()

	if err != nil {
		log.Fatalln("Something went wrong with OAuth2.0 creator: " + err.Error())
	}
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(handler.CustomPageNotFound)
	r.HandleFunc("/", handler.RootPath)
	r.HandleFunc("/login/google", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, auth.ClientConfig.AuthCodeURL("state"), http.StatusFound)
	})
	r.HandleFunc("/login/google/callback", auth.HandleCallback)

	return r
}
