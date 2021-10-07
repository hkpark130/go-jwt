package app

import (
	"golang/jwt/auth"
	"golang/jwt/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	jwtUserRepository := &repository.JwtUserRepository{DB: db}
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK,
			"text/html; charset=utf-8",
			[]byte("index"))
	})

	r.GET("/token", func(c *gin.Context) { auth.GetTokenHandler(c) })
	r.GET("/login", func(c *gin.Context) { auth.RenderLoginView(c) })
	r.POST("/api/login", func(c *gin.Context) { auth.Authentication(c) })

	r.POST("/user/register", func(c *gin.Context) { auth.RegisterHandler(c, jwtUserRepository) })
	r.GET("/user/:id", func(c *gin.Context) { auth.GetUserByIDHandler(c, jwtUserRepository) })
	r.GET("/users", func(c *gin.Context) { auth.GetUsersHandler(c, jwtUserRepository) })
	r.PUT("/user/update", func(c *gin.Context) { auth.UpdateHandler(c, jwtUserRepository) })
	r.DELETE("/user/delete", func(c *gin.Context) { auth.DeleteHandler(c, jwtUserRepository) })

	return r
}
