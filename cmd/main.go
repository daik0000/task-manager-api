package main

import (
	"log"
	"net/http"

	"github.com/daik0000/task-manager-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real env vars")
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/db-check", func(c *gin.Context) {
		sqlDB, _ := database.DB()
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"db": "unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"db": "connected"})
	})

	r.Run(":30000")
}
