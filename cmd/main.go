package main

import (
	"github.com/sayasurvey/cbm_api/router"
	"github.com/sayasurvey/cbm_api/model/database"
)

func main() {
	database.DbInit()

	router := router.GetRouter()
	router.Run(":8080")
}
