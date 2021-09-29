package repository

import "golang/jwt/domain"

type UserRepository interface {
	CreateUser(*domain.JwtUser) error
	UpdateUser(*domain.JwtUser, int) error
	GetUsers() ([]domain.JwtUser, error)
	GetUserByID(int) (domain.JwtUser, error)
	DeleteUserByID(int) error
}
