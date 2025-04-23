package main

import (
  // "github.com/soramar/CBM_api/router"
)

func main() {
  // database.DbInit()

  router := router.GetRouter()
  router.Run(":8080")
}
