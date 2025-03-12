package routes

import (
	"go_edtech_backend/controllers"
	"go_edtech_backend/middlewares"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(router *gin.Engine) {
	course := router.Group("api/v1/courses")
	{
		course.POST("/create", middlewares.AuthCheck(), controllers.Create())
		course.DELETE("/deleteCourse/:courseID", middlewares.AuthCheck(), controllers.DeleteCourse())
		course.GET("/course/:courseID", middlewares.AuthCheck(), controllers.GetCourseByID())
	}
}
