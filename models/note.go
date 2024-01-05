package models

import (
	"time"
)

type Note struct {
	ID          string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      string `gorm:"type:uuid;not null"`
	Title       string `gorm:"not null"`
	Content     string
	Shared      bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsDeleted   bool         `gorm:"type:boolean;default:false"`
	SharedNotes []SharedNote `gorm:"foreignKey:NoteID"`
}
