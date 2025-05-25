package models

import (
	"time"
	"gorm.io/gorm"
)

type Note struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Content     string         `json:"content" gorm:"type:text"`
	OwnerID     uint           `json:"owner_id" gorm:"not null"`
	Owner       User           `json:"owner"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	Version     int            `json:"version" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Collaborations []Collaboration `json:"collaborations,omitempty"`
	Operations     []Operation     `json:"operations,omitempty"`
}

type CreateNoteRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	IsPublic bool   `json:"is_public"`
}

type UpdateNoteRequest struct {
	Title    *string `json:"title,omitempty"`
	Content  *string `json:"content,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

type ShareNoteRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Permission string `json:"permission" binding:"required,oneof=read write"`
}
