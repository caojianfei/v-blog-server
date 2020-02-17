package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Nickname string
	Email string
	Password string
}
