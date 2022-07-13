package db

import (
	"errors"
)

type Show struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var shows = []Show{
	{ID: "1", Title: "Blue Train"},
	{ID: "2", Title: "Jeru"},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown"},
}

func Get(id string) (Show, error) {
	for _, a := range shows {
		if a.ID == id {
			return a, nil
		}
	}
	return Show{}, errors.New("Not found")
}

func GetAll() []Show {
	return shows
}

func Add(data Show) {
	shows = append(shows, data)
}

func Update(id string, data Show) error {
	for i, a := range shows {
		if a.ID == id {
			shows[i] = data
			return nil
		}
	}
	return errors.New("Not found")
}

func Delete(id string) error {
	idx := -1
	for i, a := range shows {
		if a.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("Not found")
	}
	shows[idx] = shows[len(shows)-1]
	shows = shows[:len(shows)-1]
	return nil
}
