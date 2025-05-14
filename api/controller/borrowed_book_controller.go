package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sayasurvey/golang/model/schema"
	"github.com/sayasurvey/golang/model/database"
	"net/http"
	"time"
)

type BorrowBookRequest struct {
	BookID        uint      `json:"book_id" binding:"required"`
	CheckoutDate  time.Time `json:"checkout_date" binding:"required"`
	ReturnDueDate time.Time `json:"return_due_date" binding:"required"`
}

type ReturnBookRequest struct {
	BorrowedBookID uint `json:"borrowed_book_id" binding:"required"`
}

func BorrowBook(c *gin.Context) {
	var request BorrowBookRequest
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

	var book schema.Book
	if err := database.Db.First(&book, request.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "本が見つかりません",
		})
		return
	}

	if !book.Loanable {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "この本は現在貸し出しできません",
		})
		return
	}

	borrowedBook := schema.BorrowedBook{
		UserID:        userID.(uint),
		BookID:        request.BookID,
		CheckoutDate:  request.CheckoutDate,
		ReturnDueDate: request.ReturnDueDate,
	}

	tx := database.Db.Begin()

	if err := tx.Create(&borrowedBook).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "貸し出し情報の保存に失敗しました",
		})
		return
	}

	if err := tx.Model(&book).Update("loanable", false).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "本の状態更新に失敗しました",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トランザクションのコミットに失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "本の貸し出しが完了しました",
		"borrowed_book": borrowedBook,
	})
}

func ReturnBook(c *gin.Context) {
	var request ReturnBookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエストボディが不正です",
		})
		return
	}

	var borrowedBook schema.BorrowedBook
	if err := database.Db.First(&borrowedBook, request.BorrowedBookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "貸し出し情報が見つかりません",
		})
		return
	}

	var book schema.Book
	if err := database.Db.First(&book, borrowedBook.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "本が見つかりません",
		})
		return
	}

	tx := database.Db.Begin()

	if err := tx.Delete(&borrowedBook).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "貸し出し情報の削除に失敗しました",
		})
		return
	}

	if err := tx.Model(&book).Update("loanable", true).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "本の状態更新に失敗しました",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "トランザクションのコミットに失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "本の返却が完了しました",
	})
}
