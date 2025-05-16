package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sayasurvey/golang/model/schema"
	"github.com/sayasurvey/golang/model/database"
	"net/http"
)

type AddToWishListRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

type WishListResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url"`
}

func AddToWishList(c *gin.Context) {
	var request AddToWishListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストボディが不正です",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証が必要です",
		})
		return
	}

	// 本の存在確認
	var book schema.Book
	if err := database.Db.First(&book, request.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "本が見つかりません",
		})
		return
	}

	// 既に追加済みかチェック
	var existingWishList schema.BorrowingWishList
	err := database.Db.Where("user_id = ? AND book_id = ?", userID, request.BookID).First(&existingWishList).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "この本は既にお気に入りに追加されています",
		})
		return
	}

	wishList := schema.BorrowingWishList{
		UserID: userID.(uint),
		BookID: request.BookID,
	}

	if err := database.Db.Create(&wishList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "お気に入りの追加に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "お気に入りに追加しました",
		"wish_list": WishListResponse{
			ID:       book.ID,
			Title:    book.Title,
			ImageUrl: book.ImageUrl,
		},
	})
}

func RemoveFromWishList(c *gin.Context) {
	bookID := c.Param("book_id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証が必要です",
		})
		return
	}

	var wishList schema.BorrowingWishList
	if err := database.Db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&wishList).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "お気に入りが見つかりません",
		})
		return
	}

	if err := database.Db.Delete(&wishList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "お気に入りの削除に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "お気に入りから削除しました",
	})
}

func GetWishList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証が必要です",
		})
		return
	}

	var wishList []schema.BorrowingWishList
	if err := database.Db.Where("user_id = ?", userID).Find(&wishList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "お気に入りリストの取得に失敗しました",
		})
		return
	}

	response := []WishListResponse{}
	if len(wishList) > 0 {
		for _, item := range wishList {
			var book schema.Book
			if err := database.Db.First(&book, item.BookID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "本の情報取得に失敗しました",
				})
				return
			}

			response = append(response, WishListResponse{
				ID:       book.ID,
				Title:    book.Title,
				ImageUrl: book.ImageUrl,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"wish_list": response,
	})
}
