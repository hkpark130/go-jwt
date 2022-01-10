package repository

import (
	"fmt"
	"golang/jwt/api/domain"

	"log"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type JwtUserRepository struct {
	DB    *gorm.DB
	Redis *redis.Client
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
	result := jwtUserRepository.DB.Where("email = ?", jwtUser.Email).First(&user)

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

func (jwtUserRepository JwtUserRepository) SetRefreshToken(email string, token string) error {
	err := jwtUserRepository.Redis.Set(email, token, 0).Err()

	if err != nil {
		log.Printf("Failed to set refresh token: %s ", err)
		return err
	}
	return nil
}

func (jwtUserRepository JwtUserRepository) GetRefreshToken(email string) (string, error) {
	val, err := jwtUserRepository.Redis.Get(email).Result()

	if err != nil {
		log.Printf("Failed to get refresh token: %s ", err)
		return "", err
	}
	return val, nil
}
