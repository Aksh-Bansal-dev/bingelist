package routes

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/Aksh-Bansal-dev/bingelist/pkg/db"
)

var FileServer = http.FileServer(http.Dir("./static"))

func Routes() {
	tmpl := template.Must(template.ParseFiles("static/index.html"))

	http.HandleFunc("/", pageHandler(tmpl))
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/movies", moviesHandler)
}

func pageHandler(tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not supported", 405)
			return
		}
		tmpl.Execute(w, db.GetAll())
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", 405)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(db.GetAll())
}

func addMovieHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	var body db.Show
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
