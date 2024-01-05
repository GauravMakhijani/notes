// internal/models/shared_note.go

package models

import (
	"time"
)

type SharedNote struct {
	ID         string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NoteID     string `gorm:"type:uuid;not null"`
	FromUserID string `gorm:"type:uuid;not null"`
	ToUserID   string `gorm:"type:uuid;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
