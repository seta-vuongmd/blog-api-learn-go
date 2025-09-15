package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags" gorm:"type:text[]"`
	UserID  uint     `json:"user_id"`
}
