package main

import (
	"github.com/sayasurvey/golang/router"
	"github.com/sayasurvey/golang/model/database"
)

func main() {
	database.DbInit()

	router := router.GetRouter()
	router.Run(":8080")
}
