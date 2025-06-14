package main

import (
	"github.com/sayasurvey/golang/model/database"
	"github.com/sayasurvey/golang/router"
	"os"
)

func main() {
	database.DbInit()

	router := router.GetRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
