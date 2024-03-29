package domain

import (
	"gorm.io/gorm"
)

type JwtUser struct {
	gorm.Model
	Email      string `binding:"required" form:"email"`
	Permission string `binding:"required" form:"permission"`
	Password   string `binding:"required" form:"password"`
	DeletedAt  gorm.DeletedAt
}

type JwtUserUseCase interface {
	CreateUser(*JwtUser) error
	UpdateUser(*JwtUser) error
	GetUsers() ([]*JwtUser, error)
	GetUserByID(int) (*JwtUser, error)
	DeleteUserByID(int) error
	CheckUser(string, string) error
	SetRefreshToken(string, string) error
	GetRefreshToken(string) (string, error)
}
