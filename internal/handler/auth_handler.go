package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/daik0000/task-manager-api/internal/auth"
	"github.com/daik0000/task-manager-api/internal/model"
)

type AuthHandler struct {
	// Có tác dụng là một struct để xử lý các yêu cầu liên quan đến xác thực người dùng trong ứng dụng web. Nó chứa một trường DB kiểu *gorm.DB, đại diện cho kết nối cơ sở dữ liệu sử dụng GORM.
	DB *gorm.DB
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (handler *AuthHandler) Register(c *gin.Context) {
	// Implementation for user registration

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// ShouldBindJSON: is a method provided by the Gin framework that binds the incoming JSON payload to the specified struct (in this case, RegisterRequest). It automatically validates the request data based on the struct tags (e.g., required, email, min=6).
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := model.User{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}
	if err := handler.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user, maybe the email is already registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": gin.H{"id": user.ID, "email": user.Email}})
}

func (handler *AuthHandler) Login(c *gin.Context) {
	// Implementation for user login

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// ShouldBindJSON: is a method provided by the Gin framework that binds the incoming JSON payload to the specified struct (in this case, LoginRequest). It automatically validates the request data based on the struct tags (e.g., required, email, min=6).
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User

	if err := handler.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !auth.CheckPassword(req.Password, user.HashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
