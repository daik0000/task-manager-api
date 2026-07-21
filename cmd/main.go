package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/daik0000/task-manager-api/internal/db"
	"github.com/daik0000/task-manager-api/internal/handler"
	"github.com/daik0000/task-manager-api/internal/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real env vars")
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	authHandler := &handler.AuthHandler{DB: database}
	taskHandler := &handler.TaskHandler{DB: database}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	tasks := r.Group("/tasks")           // Tạo một nhóm route cho các endpoint liên quan đến tasks. Nhóm này sẽ có tiền tố "/tasks" cho tất cả các route bên trong nó.
	tasks.Use(middleware.AuthRequired()) // Áp dụng middleware AuthRequired cho tất cả các route bên trong nhóm tasks. Điều này có nghĩa là bất kỳ yêu cầu nào đến các endpoint trong nhóm này sẽ phải đi qua middleware AuthRequired trước khi được xử lý bởi các handler tương ứng.
	{
		tasks.POST("", taskHandler.Create)
		tasks.GET("", taskHandler.List)
		tasks.PUT("/:id", taskHandler.Update)
		tasks.DELETE("/:id", taskHandler.Delete)
	}

	r.Run(":30000")
}
