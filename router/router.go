package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sayasurvey/golang/api/controller"
	"github.com/sayasurvey/golang/middleware"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", controller.SayHello)
	r.POST("api/login", controller.Login)

	api := r.Group("/api", middleware.JWTAuthMiddleware())
	{
		api.POST("/register", controller.Register)
		api.POST("/logout", controller.Logout)
		api.GET("/users", controller.GetUsers)
		api.GET("/books", controller.GetBooks)
		api.POST("/books", controller.CreateBook)
		api.PUT("/books/:id", controller.UpdateBook)
		api.DELETE("/books/:id", controller.DeleteBook)
		api.POST("/books/borrow", controller.BorrowBook)
		api.POST("/books/return", controller.ReturnBook)
		api.GET("/books/borrowed", controller.GetBorrowedBooks)
		api.POST("/books/wish-list", controller.AddToWishList)
		api.DELETE("/books/wish-list/:book_id", controller.RemoveFromWishList)
		api.GET("/books/wish-list", controller.GetWishList)
	}

	return r
}
