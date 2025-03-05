package routes

import (
	"go_edtech_backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("api/v1/users")
	{
		user.POST("/signup", controllers.Signup())
		user.POST("/login", controllers.Login())
	}
}
