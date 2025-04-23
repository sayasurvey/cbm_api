package router

import (
	"github.com/gin-gonic/gin"
	// "github.com/soramar/CBM_api/api/controller"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", controller.SayHello)

	return r
}
