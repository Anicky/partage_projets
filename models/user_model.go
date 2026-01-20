package models

import "time"

type User struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Email         string    `gorm:"unique" binding:"required,email"`
	Password      string    `binding:"required,min=8"`
	Comments      []Comment `gorm:"foreignKey:UserID"`
	LikedProjects []Project `gorm:"many2many:project_likes"`
}
