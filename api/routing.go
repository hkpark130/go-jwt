package app

import (
	"golang/jwt/api/handlers"
	"golang/jwt/api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	jwtUserRepository := &repository.JwtUserRepository{DB: db}
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK,
			"text/html; charset=utf-8",
			[]byte("index"))
	})

	r.GET("/token", func(c *gin.Context) { handlers.GetTokenHandler(c) })
	// r.GET("/login", func(c *gin.Context) { auth.RenderLoginView(c) })
	r.POST("/api/login", func(c *gin.Context) { handlers.Authentication(c) })

	r.POST("/user/register", func(c *gin.Context) { handlers.RegisterHandler(c, jwtUserRepository) })
	r.GET("/user/:id", func(c *gin.Context) { handlers.GetUserByIDHandler(c, jwtUserRepository) })
	r.GET("/users", func(c *gin.Context) { handlers.GetUsersHandler(c, jwtUserRepository) })
	r.PUT("/user/update", func(c *gin.Context) { handlers.UpdateHandler(c, jwtUserRepository) })
	r.DELETE("/user/delete", func(c *gin.Context) { handlers.DeleteHandler(c, jwtUserRepository) })

	return r
}
