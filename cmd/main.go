package main

import (
	"github.com/sayasurvey/golang/router"
)

func main() {
	// database.DbInit()

	router := router.GetRouter()
	router.Run(":8080")
}
