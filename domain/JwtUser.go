package domain

import (
	"time"

	"gorm.io/gorm"
)

type JwtUser struct {
	gorm.Model
	Id        int64     `gorm:"primaryKey"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"type:DATE"`
	UpdatedAt time.Time `gorm:"type:DATE"`
	DeletedAt time.Time `gorm:"type:DATE"`
}

// TODO: interface をrepositoryで実装する
// type JwtUserUseCase interface {
// 	CreateUser(*JwtUser) error
// 	UpdateUser(*JwtUser) error
// 	GetUsers() ([]*JwtUser, error)
// 	GetUserByID(int) (*JwtUser, error)
// 	DeleteUserByID(int) error
// }

type JwtUserUseCase struct {
}

func (usecase *JwtUserUseCase) CreateUser(u *JwtUser) error {
	// TODO: DB insert機能
	// err :=

	// if err != nil {
	// 	err = fmt.Errorf("[usecase.CreateUser] failed: %w ", err)
	// 	return err
	// }
	return nil
}
