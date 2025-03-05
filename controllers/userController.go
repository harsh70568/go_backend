package controllers

import (
	"go_edtech_backend/db"
	"go_edtech_backend/models"
	"go_edtech_backend/utils"
	"net/http"
	"time"

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

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Body request", "details": err.Error()})
			return
		}

		if user.Email == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is missing"})
			return
		}

		var existingUser models.User
		if err := db.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		if err := utils.VerifyPassword(user.Password, existingUser.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Password"})
			return
		}

		token, refreshToken, err := utils.GenerateNewTokens(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens", "details": err.Error()})
			return
		}

		c.SetCookie("token", token, int(48*time.Hour.Seconds()), "/", "localhost", false, true)
		c.SetCookie("refresh_token", refreshToken, int(240*time.Hour.Seconds()), "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"user":          existingUser,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}
