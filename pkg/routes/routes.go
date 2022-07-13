package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Aksh-Bansal-dev/bingelist/pkg/db"
)

type PingResponse struct {
	Message string `json:"message"`
}

var FileServer = http.FileServer(http.Dir("./static"))

func Routes() {

	http.HandleFunc("/movies", getMoviesHandler)
	http.HandleFunc("/movies", addMovieHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/ping", pingHandler)
}

var fileServer = http.FileServer(http.Dir("./static"))

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", 405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PingResponse{Message: "pong"})
}

func getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(db.GetAll())
}

func addMovieHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	var body db.Album
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Could not add movie, try again later", 400)
		return
	}
	db.Add(body)
	json.NewEncoder(w).Encode(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
}
