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

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is requitred"})
			return
		}

		var existingUser models.User
		if err := db.DB.Where("ID = ?", userID).First(&existingUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": existingUser,
		})
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var allUsers []models.User

		if err := db.DB.Find(&allUsers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying users", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users": allUsers,
		})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
			c.Abort()
			return
		}

		/* Convert email to string (since c.Get returns interface{}) */
		emailStr, ok := email.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Inavlid email format"})
			c.Abort()
			return
		}

		if res := db.DB.Model(&models.User{}).Where("email = ?", emailStr).Updates(map[string]interface{}{"token": "", "refresh_token": ""}).Error; res != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout", "details": res.Error()})
			c.Abort()
			return
		}

		c.SetCookie("token", "", -1, "/", "localhost", false, true)
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "Logout sucessfull", "success": true})
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		type change struct {
			Email string `json:"email"`
			Old   string `json:"oldPassword"`
			New   string `json:"newPassword"`
		}
		var request change
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request", "details": err.Error()})
			c.Abort()
			return
		}

		var existingUser models.User
		if err := db.DB.Where("email = ?", request.Email).First(&existingUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No user with given email"})
			return
		}

		/* Check if password stored in db is as same as of password in request */
		if err := utils.VerifyPassword(request.Old, existingUser.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password do not match"})
			return
		}

		newHashPassword, err := utils.HashPassword(request.New)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		if err := db.DB.Model(&models.User{}).Where("email = ?", request.Email).Update("password", newHashPassword).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "password updated succesfully", "success": true})
	}
}
