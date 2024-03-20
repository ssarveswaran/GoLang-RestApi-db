package main

// import "gorm.io/gorm"

type Users struct {

	Id int `gorm:"primaryKey"`
	Username string  `json:"username"`
	Password string   `json:"password"`
}
