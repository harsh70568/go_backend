package main

import (
	"go_edtech_backend/db"
	"go_edtech_backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	/* Set up database connection */
	db.ConnectDB()

	/* Set up a router */
	router := gin.Default()

	/* Setting up routes */
	routes.UserRoutes(router)

	/* Starting the server */
	err := router.Run(db.GetPort())
	if err != nil {
		log.Fatalf("Error starting the server %v", err)
	}
}
