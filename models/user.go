package models

import (
	"time"
)

type User struct {
	ID           string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt    time.Time
	Username     string       `gorm:"unique;not null"`
	PasswordHash string       `gorm:"not null"`
	Notes        []Note       `gorm:"foreignKey:UserID"`
	SharedNotes  []SharedNote `gorm:"foreignKey:ToUserID"`
}
