package repository

import (
	"fmt"
	"golang/jwt/api/domain"

	"gorm.io/gorm"
)

type JwtUserRepository struct {
	DB *gorm.DB
}

func (jwtUserRepository JwtUserRepository) CreateUser(u *domain.JwtUser) error {
	result := jwtUserRepository.DB.Create(u)

	if result.Error != nil {
		err := fmt.Errorf("[usecase.CreateUser] failed: %w ", result.Error)
		return err
	}
	return nil
}

func (jwtUserRepository JwtUserRepository) GetUserByID(i uint64) (*domain.JwtUser, error) {
	var user *domain.JwtUser
	result := jwtUserRepository.DB.First(&user, i)

	return user, result.Error
}

func (jwtUserRepository JwtUserRepository) LoginEmailPassword(jwtUser *domain.JwtUser) (*domain.JwtUser, error) {
	var user *domain.JwtUser
	result := jwtUserRepository.DB.Where("email = ? AND password = ?", jwtUser.Email, jwtUser.Password).First(&user)

	return user, result.Error
}

func (jwtUserRepository JwtUserRepository) GetUsers() ([]*domain.JwtUser, error) {
	var users []*domain.JwtUser
	result := jwtUserRepository.DB.Find(&users)

	return users, result.Error
}

func (jwtUserRepository JwtUserRepository) UpdateUser(i uint64, u *domain.JwtUser) (*domain.JwtUser, error) {
	var user *domain.JwtUser
	result := jwtUserRepository.DB.First(&user, i).Updates(u)

	return user, result.Error
}

func (jwtUserRepository JwtUserRepository) DeleteUserByID(i uint64) error {
	var user *domain.JwtUser
	result := jwtUserRepository.DB.Delete(&user, i)

	return result.Error
}