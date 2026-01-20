package models

import "time"

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ProjectID uint `json:"project_id"`
	UserID    uint
	Content   string
}
