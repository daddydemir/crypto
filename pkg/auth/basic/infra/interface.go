package infra

import (
	"fmt"

	"github.com/daddydemir/crypto/pkg/auth/basic/domain"
	"github.com/daddydemir/crypto/pkg/token/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	database     *gorm.DB
	tokenService *jwt.TokenService
}

func NewRepository(db *gorm.DB, tokenService *jwt.TokenService) *Repository {
	return &Repository{
		database:     db,
		tokenService: tokenService,
	}
}

func (r *Repository) CreateUser(user domain.User) error {
	user.ID = uuid.New().String()
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return r.database.Create(&user).Error
}

func (r *Repository) Login(username, password string) (*domain.LoginResponse, error) {
	user, err := r.findUserByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	token, err := r.tokenService.GenerateToken(username)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{Token: token, Username: user.Username}, nil
}

func (r *Repository) findUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.database.Where("username = ? and is_active = ?", username, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
