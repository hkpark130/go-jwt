package app

import (
	"golang/jwt/api/handlers"
	"golang/jwt/api/middleware"
	"golang/jwt/api/repository"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client) *gin.Engine {
	jwtUserRepository := &repository.JwtUserRepository{DB: db, Redis: redis}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:8300",
			"https://localhost:8300",
			"http://hkpark130.p-e.kr:8300",
			"https://hkpark130.p-e.kr:8300",
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
	r.POST("/register", func(c *gin.Context) { handlers.RegisterHandler(c, jwtUserRepository) })

	// user API router
	r.Group("/user", middleware.Authorization(jwtUserRepository)).
		GET("/token", func(c *gin.Context) { handlers.GetTokenHandler(c, jwtUserRepository) }).
		GET("/logout", func(c *gin.Context) { handlers.LogoutHandler(c, jwtUserRepository) }).
		GET("/admin", func(c *gin.Context) { handlers.AdminHandler(c, jwtUserRepository) }).
		GET("/:id", func(c *gin.Context) { handlers.GetUserByIDHandler(c, jwtUserRepository) }).
		GET("/users", func(c *gin.Context) { handlers.GetUsersHandler(c, jwtUserRepository) }).
		PUT("/update", func(c *gin.Context) { handlers.UpdateHandler(c, jwtUserRepository) }).
		DELETE("/delete", func(c *gin.Context) { handlers.DeleteHandler(c, jwtUserRepository) })

	return r
}
