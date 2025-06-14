package schema

import (
	"gorm.io/gorm"
	"time"
)

type Role string

const (
	UserRole Role = "USER"
	AdminRole Role = "ADMIN"
)

type User struct {
	gorm.Model
	Name              	string `gorm:"type:varchar(255);not null"                         validate:"required"`
	Email             	string `gorm:"type:varchar(255);uniqueIndex;not null"             validate:"required,email"`
	Password          	string `gorm:"type:varchar(255);not null"                         validate:"required,min=8"`
	Role              	Role   `gorm:"type:varchar(10);default:'USER';not null" validate:"required"`
	Books 				[]Book
	BorrowedBooks 		[]BorrowedBook
	BorrowingWishLists 	[]BorrowingWishList
}

type Book struct {
	gorm.Model
	UserId		uint	 `validate:"required"`
	User      User
	Title  		string `gorm:"type:varchar(255);not null" validate:"required"`
	ImageUrl  string `gorm:"type:varchar(255);not null" validate:"required"`
	Loanable  bool   `gorm:"not null"                   validate:"required"`
}

type BorrowedBook struct {
	gorm.Model
	UserID        uint      `gorm:"not null" validate:"required"`
	BookID        uint      `gorm:"not null" validate:"required"`
	CheckoutDate  time.Time `gorm:"not null" validate:"required"`
	ReturnDueDate time.Time `gorm:"not null" validate:"required"`
}

type BorrowingWishList struct {
	gorm.Model
	UserID    uint      `gorm:"not null" validate:"required"`
	BookID    uint      `gorm:"not null" validate:"required"`
}

type InvalidatedToken struct {
	gorm.Model
	Token     string    `gorm:"primaryKey" validate:"required"`
	ExpiresAt time.Time `gorm:"not null" validate:"required"`
}

type LoginRequest struct {
	Email    string 		`json:"email"    validate:"required"`
	Password string 		`json:"password" validate:"required"`
}
