package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Jwt struct {
	AccessToken  string
	RefreshToken string
}

type JwtUser struct {
	gorm.Model
	Id        int64     `gorm:"primaryKey"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"type:DATE"`
	UpdatedAt time.Time `gorm:"type:DATE"`
	DeletedAt time.Time `gorm:"type:DATE"`
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

	result := db.Create(&JwtUser{Email: "t@test.com", Password: "t"})

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
	user := &JwtUser{Email: email, Password: password}

	c.JSON(http.StatusOK, user)
}
