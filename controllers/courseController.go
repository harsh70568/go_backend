package controllers

import (
	"go_edtech_backend/db"
	"go_edtech_backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
			return
		}

		var course models.Course
		if err := c.ShouldBindJSON(&course); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Inavlid request body", "details": err.Error()})
			return
		}

		var user models.User
		if err := db.DB.Where("email = ?", email).Find(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found", "details": err.Error()})
			return
		}

		course.Owner = user.ID
		course.Students = 0
		course.Ratings = 0
		course.CreatedAt = time.Now()
		course.UpdatedAt = time.Now()

		if err := db.DB.Create(&course).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in inserting to database", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Course created succesfully", "id": course.ID})
	}
}

func GetCourseByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("courseID")
		if courseID == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Course ID is not valid"})
			c.Abort()
			return
		}

		var course models.Course
		if err := db.DB.Where("id = ?", courseID).Find(&course).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Course ID not found"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"course": course})
	}
}

func DeleteCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("courseID")
		if courseID == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "course ID is not valid"})
			c.Abort()
			return
		}

		id, err := strconv.ParseUint(courseID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		if err := db.DB.Delete(&models.Course{}, uint(id)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete course"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "course deleted succesfully"})
	}
}
