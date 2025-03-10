package models

import "time"

type Course struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Owner       uint      `json:"owner" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null;size:50"; validate:"required;min=2,max=50"`
	Description string    `json:"description" gorm:"not null;size:200"; validate:"required;min=30,max=200"`
	Tag         uint      `json:"tag" gorm:"not null"`
	Price       int       `json:"price" gorm:"not null"`
	Students    int       `json:"students" gorm:"default:0"`
	Ratings     int       `json:"ratings", gorm:"default:0"`
	CreatedAt   time.Time `json:"createdAt", gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt", gorm:"autoUpdateTime"`
}
