package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Authentication Middleware
func SimpleAuth() gin.HandlerFunc {
	auth := os.Getenv("AUTH_HEADER")
	if auth == "" {
		panic("AUTH_HEADER environment variable not set")
	}
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != auth {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Modpack{}, &LatestModpack{})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://modpack-manager.octsrv.org", "https://localhost:5713"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	authGroup := r.Group("/").Use(SimpleAuth())
	{
		// Modpack CRUD operations
		authGroup.GET("/modpacks", GetModpacks)
		authGroup.POST("/modpacks", CreateModpack)
		authGroup.GET("/modpacks/:id", GetModpack)
		authGroup.PATCH("/modpacks/:id", UpdateModpack)
		authGroup.DELETE("/modpacks/:id", DeleteModpack)
		// Publish modpacks
		authGroup.GET("/published", GetLatestModpacks)
		authGroup.PUT("/publish/:server/modpack/:modpack_id", SetLatestModpack)
		authGroup.DELETE("/published/:server", DeleteLatestModpack)
	}
	r.GET("/published/:server", GetLatestModpack)

	// Support for legacy URL
	r.GET("/servers/:server/modpack", GetLatestModpack)

	r.Run(":8080")
}
