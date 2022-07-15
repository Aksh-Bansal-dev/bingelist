package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/Aksh-Bansal-dev/bingelist/pkg/db"
)

type ShowRes struct {
	ID        string
	Title     string
	Upvotes   int
	CanUpvote bool
}

var FileServer = http.FileServer(http.Dir("./static"))

func Routes() {
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	setup()

	http.HandleFunc("/", pageHandler(tmpl))
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/add", moviesHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/login/redirect", googleCallbackHandler)
}

func pageHandler(tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not supported", 405)
			return
		}
		// separate api for getting if upvoted a show or not
		res := []ShowRes{}
		shows, upvotes := db.GetAll()
		for _, show := range shows {
			votes := 0
			canUpvote := false
			for _, upvote := range upvotes {
				if upvote.ShowID == show.ID {
					votes++
				}
			}
			res = append(res, ShowRes{ID: show.ID, Title: show.Title, Upvotes: votes, CanUpvote: canUpvote})
		}
		tmpl.Execute(w, res)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
	var body db.Show
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Could not add movie, try again later", 400)
		return
	}
	if body.Title == "" {
		http.Error(w, "Title must be non-empty", 400)
		return
	}
	body.ID = fmt.Sprint(rand.Int())
	db.AddShow(body)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
	token := r.Header.Get("Authorization")
	if !db.DoesUserExist(token) {
		http.Error(w, "User doesnt exist", http.StatusBadRequest)
		return
	}
	var body db.Upvote
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Try again later", 400)
		return
	}

	err = db.AddVote(body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
}
