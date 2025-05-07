package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sayasurvey/golang/api/controller"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", controller.SayHello)
	api := r.Group("/api")
	{
		api.GET("/books", controller.GetBooks)
		api.POST("/books", controller.CreateBook)
	}

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}

	return r
}
