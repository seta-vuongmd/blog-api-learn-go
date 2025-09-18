package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Tags    pq.StringArray `json:"tags" gorm:"type:text[]"`
	UserID  uint           `json:"user_id"`
}
