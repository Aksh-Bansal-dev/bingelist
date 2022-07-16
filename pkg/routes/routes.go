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
	log.Println("/ping")
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		log.Println("Method not supported")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func moviesHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	log.Println("/add")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed.")
		return
	}
	token := r.Header.Get("Authorization")
	if v, err := db.DoesUserExist(database, token); err != nil || !v {
		if err != nil {
			log.Println(err)
			http.Error(w, "Something wend wrong with DB", 400)
			return
		} else {
			log.Println("User does not exists")
		}
		http.Error(w, "User does not exists", http.StatusBadRequest)
		return
	}
	var body db.Show
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Could not add movie, try again later", 400)
		log.Println(err)
		return
	}
	if body.Title == "" {
		http.Error(w, "Title must be non-empty", 400)
		log.Println("Title must be non-empty")
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
	log.Println("/init-data")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed.")
		return
	}
	token := r.Header.Get("Authorization")
	if v, err := db.DoesUserExist(database, token); err != nil || !v {
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong in DB", 400)
			return
		} else {
			log.Println("User does not exists")
		}
		token = ""
	}
	res := []InitDataRow{}
	data, err := db.GetShows(database)
	if err != nil {
		http.Error(w, "Something went wrong in DB", 400)
		log.Println(err)
		return
	}
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
	log.Println("/vote")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed.")
		return
	}
	token := r.Header.Get("Authorization")
	if v, err := db.DoesUserExist(database, token); err != nil || !v {
		if err != nil {
			log.Println(err)
			http.Error(w, "Something wend wrong with DB", 400)
			return
		} else {
			log.Println("User does not exists")
		}
		http.Error(w, "User does not exists", http.StatusBadRequest)
		return
	}
	var body VoteBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Try again later", 400)
		fmt.Println(err)
		return
	}
	showId, err := strconv.ParseUint(body.ShowID, 10, 0)
	if err != nil {
		http.Error(w, "Try again later", 400)
		fmt.Println(err)
		return
	}
	vote := db.Upvote{
		ShowID: uint(showId),
		UserID: body.UserID,
	}
	err = db.AddVote(database, vote)
	if err != nil {
		http.Error(w, err.Error(), 400)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
}
