package auth

import (
	"golang/jwt/domain"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAuthenticationHandler(t *testing.T) {
	dsn := "host=db user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Errorf("Failed to connect to database. %s", err)
	}

	result := db.Create(&JwtUser{Email: "t@test.com", Password: "t"})
	if result.Error != nil {
		t.Errorf("Failed to insert. %s ", result.Error)
	}

	var user domain.JwtUser
	db.Where("Email = ?", "t@test.com").Find(&user)
	if e := user.Email; e != "t@test.com" {
		t.Errorf("Email doesn't match. %s != %s", "t@test.com", e)
	}

	del_result := db.Delete(&user)
	if del_result.Error != nil {
		t.Errorf("Failed to delete. %s ", result.Error)
	}
}
