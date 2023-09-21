package models

import "gorm.io/gorm"

// Post table schema
type Post struct {
	gorm.Model
	Title string
	Body string
}