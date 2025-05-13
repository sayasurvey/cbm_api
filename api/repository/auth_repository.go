package repository

import (
	"github.com/sayasurvey/golang/model/schema"
	"github.com/sayasurvey/golang/model/database"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) CreateUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := schema.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     schema.UserRole,
	}

	result := database.Db.Create(&user)
	return result.Error
}

func (r *AuthRepository) FindUserByEmail(email string) (*schema.User, error) {
	var user schema.User
	err := database.Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) ValidatePassword(user *schema.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (r *AuthRepository) GenerateToken(user *schema.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func (r *AuthRepository) GetAllUsers() ([]schema.User, error) {
	var users []schema.User
	err := database.Db.Find(&users).Error
	return users, err
}
