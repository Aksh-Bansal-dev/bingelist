package db

import (
	"errors"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist Artist  `json:"artist"`
	Price  float64 `json:"price"`
}

type Artist struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: Artist{Name: "John Coltrane", Age: 56}, Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: Artist{Name: "Gerry Mulligan", Age: 43}, Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: Artist{Name: "Sarah Vaughan", Age: 32}, Price: 39.99},
}

func Get(id string) (Album, error) {
	for _, a := range albums {
		if a.ID == id {
			return a, nil
		}
	}
	return Album{}, errors.New("Not found")
}

func GetAll() []Album {
	return albums
}

func Add(data Album) {
	albums = append(albums, data)
}

func Update(id string, data Album) error {
	for i, a := range albums {
		if a.ID == id {
			albums[i] = data
			return nil
		}
	}
	return errors.New("Not found")
}

func Delete(id string) error {
	idx := -1
	for i, a := range albums {
		if a.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("Not found")
	}
	albums[idx] = albums[len(albums)-1]
	albums = albums[:len(albums)-1]
	return nil
}
