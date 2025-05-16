package repository

import (
	"github.com/sayasurvey/golang/model/database"
	"github.com/sayasurvey/golang/model/schema"
)

type BorrowingWishListRepository struct{}

func NewBorrowingWishListRepository() *BorrowingWishListRepository {
	return &BorrowingWishListRepository{}
}

func (r *BorrowingWishListRepository) FindBookByID(bookID uint) (*schema.Book, error) {
	var book schema.Book
	if err := database.Db.First(&book, bookID).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BorrowingWishListRepository) FindWishListByUserIDAndBookID(userID, bookID uint) (*schema.BorrowingWishList, error) {
	var wishList schema.BorrowingWishList
	if err := database.Db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&wishList).Error; err != nil {
		return nil, err
	}
	return &wishList, nil
}

func (r *BorrowingWishListRepository) CreateWishList(userID, bookID uint) (*schema.BorrowingWishList, error) {
	wishList := schema.BorrowingWishList{
		UserID: userID,
		BookID: bookID,
	}
	if err := database.Db.Create(&wishList).Error; err != nil {
		return nil, err
	}
	return &wishList, nil
}

func (r *BorrowingWishListRepository) DeleteWishList(wishList *schema.BorrowingWishList) error {
	return database.Db.Delete(wishList).Error
}

func (r *BorrowingWishListRepository) GetWishListByUserID(userID uint) ([]schema.BorrowingWishList, error) {
	var wishList []schema.BorrowingWishList
	if err := database.Db.Where("user_id = ?", userID).Find(&wishList).Error; err != nil {
		return nil, err
	}
	return wishList, nil
}
