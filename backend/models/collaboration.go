package models

import (
	"time"
	"gorm.io/gorm"
)

type Collaboration struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	NoteID     uint           `json:"note_id" gorm:"not null"`
	Note       Note           `json:"note"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	User       User           `json:"user"`
	Permission string         `json:"permission" gorm:"not null;check:permission IN ('read', 'write')"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// Operation represents a real-time collaborative operation
type Operation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	NoteID    uint      `json:"note_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user"`
	Type      string    `json:"type" gorm:"not null"` // insert, delete, retain
	Position  int       `json:"position"`
	Content   string    `json:"content"`
	Length    int       `json:"length"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}
