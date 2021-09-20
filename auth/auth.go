package auth

import (
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

type User struct {
	gorm.Model
	Email    string
	Password string
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

	db.Create(&User{Email: "t@test.com", Password: "t"})

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
	user := &User{Email: email, Password: password}

	c.JSON(http.StatusOK, user)
}
