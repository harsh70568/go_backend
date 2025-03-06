package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName    string    `json:"first_name" validate:"required,min=2,max=20"`
	LastName     string    `json:"last_name" validate:"required,min=2,max=20"`
	Email        string    `json:"email" validate:"email,required" gorm:"unique"`
	Password     string    `json:"password" validate:"required,min=8,max=16"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
