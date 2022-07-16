package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Aksh-Bansal-dev/bingelist/pkg/db"
	"gorm.io/gorm"
)

type InitDataRow struct {
	ID        string
	Title     string
	Upvotes   int
	CanUpvote bool
}

type VoteBody struct {
	ShowID string `json:"showId"`
	UserID string `json:"userId"`
}

var FileServer = http.FileServer(http.Dir("./static"))

func Routes(database *gorm.DB) {
	static := http.Dir("./static")
	setup()

	http.Handle("/", http.FileServer(static))

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		moviesHandler(w, r, database)
	})
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		voteHandler(w, r, database)
	})
	http.HandleFunc("/init-data", func(w http.ResponseWriter, r *http.Request) {
		initDataHandler(w, r, database)
	})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/login/redirect", func(w http.ResponseWriter, r *http.Request) {
		googleCallbackHandler(w, r, database)
	})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func moviesHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
	log.Println("/add")
	token := r.Header.Get("Authorization")
	if !db.DoesUserExist(database, token) {
		http.Error(w, "User doesnt exist", http.StatusBadRequest)
		return
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
	db.AddShow(database, &body)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
}

func initDataHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
	log.Println("/init-data")
	token := r.Header.Get("Authorization")
	if !db.DoesUserExist(database, token) {
		token = ""
	}
	res := []InitDataRow{}
	data := db.GetShows(database)
	for _, show := range data {
		canUpvote := true
		if token == "" {
			canUpvote = false
		}
		for _, upvote := range show.Upvotes {
			if upvote.ShowID == show.ID && upvote.UserID == token {
				canUpvote = false
			}
		}
		res = append(res, InitDataRow{
			ID:        fmt.Sprint(show.ID),
			Title:     show.Title,
			Upvotes:   len(show.Upvotes),
			CanUpvote: canUpvote,
		})
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func voteHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
	log.Println("/vote")
	token := r.Header.Get("Authorization")
	if !db.DoesUserExist(database, token) {
		http.Error(w, "User doesnt exist", http.StatusBadRequest)
		return
	}
	var body VoteBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Try again later", 400)
		return
	}
	showId, err := strconv.ParseUint(body.ShowID, 10, 0)
	if err != nil {
		http.Error(w, "Try again later", 400)
		return
	}
	vote := db.Upvote{
		ShowID: uint(showId),
		UserID: body.UserID,
	}
	err = db.AddVote(database, vote)
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
