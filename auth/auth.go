package auth

import (
	"golang/jwt/domain"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Jwt struct {
	AccessToken  string
	RefreshToken string
}

func GetTokenHandler(c *gin.Context) {
	dsn := "host=db user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	output := "success"

	if err != nil {
		output = "fail"
	}

	result := db.Create(&domain.JwtUser{Email: "t@test.com", Password: "t"})

	if result.Error != nil {
		output = "fail"
	}

	c.Data(http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(output))
}

func RenderLoginView(c *gin.Context) {
	c.HTML(http.StatusOK,
		"login.html",
		gin.H{})
}

func Authentication(c *gin.Context) {
	email := c.DefaultPostForm("email", "")
	password := c.DefaultPostForm("password", "")
	user := &domain.JwtUser{Email: email, Password: password}

	c.JSON(http.StatusOK, user)
}
