package routes

import (
	"go_edtech_backend/controllers"
	"go_edtech_backend/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("api/v1/users")
	{
		user.POST("/signup", controllers.Signup())
		user.POST("/login", controllers.Login())
		user.POST("/logout", middlewares.AuthCheck(), controllers.Logout())
		user.POST("changePassword", middlewares.AuthCheck(), controllers.ChangePassword())
		user.GET("/getUser/:userID", controllers.GetUser())
		user.GET("/getAllUsers", controllers.GetAllUsers())
	}
}
