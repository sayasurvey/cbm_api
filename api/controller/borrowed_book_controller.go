package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sayasurvey/golang/api/repository"
	"net/http"
	"time"
)

type BorrowBookRequest struct {
	BookID        uint   `json:"book_id" binding:"required"`
	CheckoutDate  string `json:"checkout_date" binding:"required"`
	ReturnDueDate string `json:"return_due_date" binding:"required"`
}

type ReturnBookRequest struct {
	BorrowedBookID uint `json:"borrowed_book_id" binding:"required"`
}

type BorrowedBookResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	ImageUrl      string `json:"image_url"`
	CheckoutDate  string `json:"checkout_date"`
	ReturnDueDate string `json:"return_due_date"`
}

var borrowedBookRepo = repository.NewBorrowedBookRepository()

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

	book, err := borrowedBookRepo.FindBookByID(request.BookID)
	if err != nil {
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

	checkoutDate, err := time.Parse("2006-01-02", request.CheckoutDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "貸出日の形式が正しくありません。YYYY-MM-DD形式で入力してください",
		})
		return
	}

	returnDueDate, err := time.Parse("2006-01-02", request.ReturnDueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "返却予定日の形式が正しくありません。YYYY-MM-DD形式で入力してください",
		})
		return
	}

	borrowedBook, err := borrowedBookRepo.CreateBorrowedBook(userID.(uint), request.BookID, checkoutDate, returnDueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "貸し出し処理に失敗しました",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "本の貸し出しが完了しました",
		"borrowed_book": gin.H{
			"id":              borrowedBook.ID,
			"user_id":         borrowedBook.UserID,
			"book_id":         borrowedBook.BookID,
			"checkout_date":   request.CheckoutDate,
			"return_due_date": request.ReturnDueDate,
		},
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

	borrowedBook, err := borrowedBookRepo.FindBorrowedBookByID(request.BorrowedBookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "貸し出し情報が見つかりません",
		})
		return
	}

	if err := borrowedBookRepo.ReturnBook(borrowedBook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "返却処理に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "本の返却が完了しました",
	})
}

func GetBorrowedBooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "認証が必要です",
		})
		return
	}

	borrowedBooks, err := borrowedBookRepo.GetBorrowedBooksByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "貸出情報の取得に失敗しました",
		})
		return
	}

	response := []BorrowedBookResponse{}
	if len(borrowedBooks) > 0 {
		for _, borrowedBook := range borrowedBooks {
			book, err := borrowedBookRepo.FindBookByID(borrowedBook.BookID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "本の情報取得に失敗しました",
				})
				return
			}

			response = append(response, BorrowedBookResponse{
				ID:            borrowedBook.ID,
				Title:         book.Title,
				ImageUrl:      book.ImageUrl,
				CheckoutDate:  borrowedBook.CheckoutDate.Format("2006-01-02"),
				ReturnDueDate: borrowedBook.ReturnDueDate.Format("2006-01-02"),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"borrowed_books": response,
	})
}
