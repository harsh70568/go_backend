package controllers

import (
	"go_edtech_backend/db"
	"go_edtech_backend/models"
	"go_edtech_backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Inavlid request body", "detail": err.Error()})
			return
		}

		/* Check if user already exists */
		var existingUser models.User
		if err := db.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		/* Hash the password */
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}
		user.Password = hashedPassword

		/* Create user in database */
		if err := db.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user", "detail": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": "true", "id": user.ID})
	}
}
