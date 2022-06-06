package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Password  string
	Name      string
	Telephone string
}
