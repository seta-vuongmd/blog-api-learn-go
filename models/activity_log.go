package models

import "gorm.io/gorm"

type ActivityLog struct {
	gorm.Model
	Action   string
	PostID   uint
	LoggedAt int64
}
