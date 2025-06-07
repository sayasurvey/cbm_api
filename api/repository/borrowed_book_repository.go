package repository

import (
	"github.com/sayasurvey/cbm_api/model/database"
	"github.com/sayasurvey/cbm_api/model/schema"
	"time"
)

type BorrowedBookRepository struct{}

func NewBorrowedBookRepository() *BorrowedBookRepository {
	return &BorrowedBookRepository{}
}

func (r *BorrowedBookRepository) FindBookByID(bookID uint) (*schema.Book, error) {
	var book schema.Book
	err := database.Db.First(&book, bookID).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BorrowedBookRepository) CreateBorrowedBook(userID uint, bookID uint, checkoutDate, returnDueDate time.Time) (*schema.BorrowedBook, error) {
	borrowedBook := schema.BorrowedBook{
		UserID:        userID,
		BookID:        bookID,
		CheckoutDate:  checkoutDate,
		ReturnDueDate: returnDueDate,
	}

	tx := database.Db.Begin()

	if err := tx.Create(&borrowedBook).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&schema.Book{}).Where("id = ?", bookID).Update("loanable", false).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &borrowedBook, nil
}

func (r *BorrowedBookRepository) FindBorrowedBookByID(id uint) (*schema.BorrowedBook, error) {
	var borrowedBook schema.BorrowedBook
	err := database.Db.First(&borrowedBook, id).Error
	if err != nil {
		return nil, err
	}
	return &borrowedBook, nil
}

func (r *BorrowedBookRepository) ReturnBook(borrowedBook *schema.BorrowedBook) error {
	tx := database.Db.Begin()

	if err := tx.Delete(borrowedBook).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&schema.Book{}).Where("id = ?", borrowedBook.BookID).Update("loanable", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *BorrowedBookRepository) GetBorrowedBooksByUserID(userID uint) ([]schema.BorrowedBook, error) {
	var borrowedBooks []schema.BorrowedBook
	err := database.Db.Where("user_id = ?", userID).Find(&borrowedBooks).Error
	if err != nil {
		return nil, err
	}
	return borrowedBooks, nil
}
