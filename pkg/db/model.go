package db

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type User struct {
	Email   string
	UserID  string   `gorm:"primaryKey"`
	Upvotes []Upvote `gorm:"foreignkey:UserID;references:UserID"`
}

type Show struct {
	ID      uint `gorm:"primaryKey"`
	Title   string
	Upvotes []Upvote `gorm:"foreignkey:ShowID;references:ID"`
}

type Upvote struct {
	ID     uint   `gorm:"primaryKey"`
	ShowID uint   `json:"showId"`
	UserID string `json:"userId"`
}

func GetShows(db *gorm.DB) []Show {
	var res []Show
	err := db.Model(&Show{}).Preload("Upvotes").Find(&res).Error
	if err != nil {
		log.Println(err)
	}
	return res
}

func AddShow(db *gorm.DB, data *Show) {
	db.Create(data)
}

func AddVote(db *gorm.DB, data Upvote) error {
	var rows []Upvote
	err := db.Where("show_id = ? AND user_id = ?", data.ShowID, data.UserID).Find(&rows).Error
	if err != nil {
		log.Println(err)
	}
	if len(rows) > 0 {
		return errors.New("Cannot vote twice")
	}
	db.Create(&data)
	return nil
}

func AddUser(db *gorm.DB, email string) string {
	hash := encrypt(email)
	data := User{Email: email, UserID: hash}
	var rows []User
	err := db.Where("email = ?", email).Find(&rows).Error
	if err != nil {
		log.Println(err)
	}
	if len(rows) == 0 {
		db.Create(&data)
	}
	return hash
}

func DoesUserExist(db *gorm.DB, token string) bool {
	var rows []User
	err := db.Where("user_id = ?", token).Find(&rows).Error
	if err != nil {
		log.Println(err)
	}
	return len(rows) > 0
}
