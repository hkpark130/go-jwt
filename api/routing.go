package app

import (
	"golang/jwt/api/handlers"
	"golang/jwt/api/middleware"
	"golang/jwt/api/repository"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	jwtUserRepository := &repository.JwtUserRepository{DB: db}
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost",
			"https://localhost",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"PUT",
			"DELETE",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	r.POST("/api/login", middleware.LoginFormValidation(), func(c *gin.Context) { handlers.Login(c, jwtUserRepository) })

	// user API router
	r.Group("/user", middleware.Authorization()).
		GET("/token", func(c *gin.Context) { handlers.GetTokenHandler(c) }).
		POST("/register", func(c *gin.Context) { handlers.RegisterHandler(c, jwtUserRepository) }).
		GET("/:id", func(c *gin.Context) { handlers.GetUserByIDHandler(c, jwtUserRepository) }).
		GET("/users", func(c *gin.Context) { handlers.GetUsersHandler(c, jwtUserRepository) }).
		PUT("/update", func(c *gin.Context) { handlers.UpdateHandler(c, jwtUserRepository) }).
		DELETE("/delete", func(c *gin.Context) { handlers.DeleteHandler(c, jwtUserRepository) })

	return r
}
