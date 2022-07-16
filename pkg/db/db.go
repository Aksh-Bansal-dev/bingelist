package db

import (
	"errors"
)

type User struct {
	Email  string `json:"email"`
	UserId string `json:"userId"`
}

type Show struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
type Upvote struct {
	ShowID string `json:"showId"`
	UserID string `json:"userId"`
}

var shows = []Show{
	{ID: "1", Title: "Blue Train"},
	{ID: "2", Title: "Jeru"},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown"},
}

var upvotes = []Upvote{
	{ShowID: "1", UserID: "001"},
	{ShowID: "1", UserID: "002"},
	{ShowID: "1", UserID: "deb83252b4353d79bf3fa48bb81df35037426dfd0e6fe7502721a4781038694c"},
	{ShowID: "1", UserID: "004"},
	{ShowID: "2", UserID: "001"},
}

var users = []User{
	{Email: "randomguy@gmail.com", UserId: "deb83252b4353d79bf3fa48bb81df35037426dfd0e6fe7502721a4781038694c"},
}

func Get(id string) (Show, error) {
	for _, a := range shows {
		if a.ID == id {
			return a, nil
		}
	}
	return Show{}, errors.New("Not found")
}

func GetAll() ([]Show, []Upvote) {
	return shows, upvotes
}

func AddShow(data Show) {
	shows = append(shows, data)
}

func AddVote(data Upvote) error {
	for _, vote := range upvotes {
		if vote == data {
			return errors.New("Cannot vote twice")
		}
	}
	upvotes = append(upvotes, data)
	return nil
}

func AddUser(email string) string {
	hash := encrypt(email)
	data := User{Email: email, UserId: hash}
	for _, u := range users {
		if u.Email == email {
			return hash
		}
	}
	users = append(users, data)
	return hash
}

func DoesUserExist(token string) bool {
	for _, u := range users {
		if u.UserId == token {
			return true
		}
	}
	return false
}
