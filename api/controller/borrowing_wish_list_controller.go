package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sayasurvey/golang/api/repository"
	"net/http"
	"strconv"
	"github.com/sayasurvey/golang/api/helper"
)

type AddToWishListRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

type WishListResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"imageUrl"`
}

type WishListResponseWrapper struct {
	WishList    []WishListResponse `json:"wishList"`
	CurrentPage int                `json:"currentPage"`
	LastPage    int                `json:"lastPage"`
	PerPage     int                `json:"perPage"`
}

var wishListRepo = repository.NewBorrowingWishListRepository()

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

	book, err := wishListRepo.FindBookByID(request.BookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "本が見つかりません",
		})
		return
	}

	_, err = wishListRepo.FindWishListByUserIDAndBookID(userID.(uint), request.BookID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "この本は既にお気に入りに追加されています",
		})
		return
	}

	_, err = wishListRepo.CreateWishList(userID.(uint), request.BookID)
	if err != nil {
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
	bookIDStr := c.Param("book_id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不正な本IDです",
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

	wishList, err := wishListRepo.FindWishListByUserIDAndBookID(userID.(uint), uint(bookID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "お気に入りが見つかりません",
		})
		return
	}

	if err := wishListRepo.DeleteWishList(wishList); err != nil {
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
	page := 1
	perPage := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if perPageStr := c.Query("perPage"); perPageStr != "" {
		if parsedPerPage, err := strconv.Atoi(perPageStr); err == nil && parsedPerPage > 0 {
			perPage = parsedPerPage
		}
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証が必要です",
		})
		return
	}

	wishList, err := wishListRepo.GetWishListByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "お気に入りリストの取得に失敗しました",
		})
		return
	}

	response := []WishListResponse{}
	if len(wishList) > 0 {
		for _, item := range wishList {
			book, err := wishListRepo.FindBookByID(item.BookID)
			if err != nil {
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

	paginatedWishList, currentPage, lastPage := helper.Pagination(response, page, perPage)

	c.JSON(http.StatusOK, WishListResponseWrapper{
		WishList:    paginatedWishList,
		CurrentPage: currentPage,
		LastPage:    lastPage,
		PerPage:     perPage,
	})
}
