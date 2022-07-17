package db

import (
	"errors"
	"os"

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

func GetShows(db *gorm.DB) ([]Show, error) {
	var res []Show
	err := db.Model(&Show{}).Preload("Upvotes").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func AddShow(db *gorm.DB, data *Show) {
	db.Create(data)
}

func DeleteShow(db *gorm.DB, showId uint) error {
	err := db.Delete(&Show{}, showId).Error
	if err != nil {
		return err
	}
	return nil
}

func AddVote(db *gorm.DB, data Upvote) error {
	var rows []Upvote
	err := db.Where("show_id = ? AND user_id = ?", data.ShowID, data.UserID).Find(&rows).Error
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		return errors.New("Cannot vote twice")
	}
	db.Create(&data)
	return nil
}

func AddUser(db *gorm.DB, email string) (string, error) {
	hash := encrypt(email)
	data := User{Email: email, UserID: hash}
	var rows []User
	err := db.Where("email = ?", email).Find(&rows).Error
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		db.Create(&data)
	}
	return hash, nil
}

func DoesUserExist(db *gorm.DB, token string) (int, error) {
	var rows []User
	err := db.Where("user_id = ?", token).Find(&rows).Error
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	} else if rows[0].Email == os.Getenv("ADMIN") {
		return 2, nil
	}
	return 1, nil
}
