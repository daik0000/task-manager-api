package handler

import (
	"net/http"
	"strconv"

	"github.com/daik0000/task-manager-api/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

type TaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (handler *TaskHandler) Create(c *gin.Context) {
	var req TaskRequest
	userID := c.MustGet("userID").(uint) // MustGet: is a method provided by the Gin framework that retrieves a value from the context. In this case, it retrieves the userID that was set in the AuthRequired middleware. The value is then type-asserted to uint.

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      model.TaskStatus(req.Status),
		OwnerID:     userID,
	}

	if task.Status == "" {
		task.Status = model.StatusTodo
	}

	if err := handler.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": task})
}

func (handler *TaskHandler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var tasks []model.Task
	handler.DB.Where("owner_id = ?", userID).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (handler *TaskHandler) getOwnerTask(c *gin.Context) (*model.Task, bool) {
	userID := c.MustGet("userID").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return nil, false
	}

	var task model.Task
	if err := handler.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return nil, false
	}

	if task.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this task"})
		return nil, false
	}

	return &task, true
}

func (handler *TaskHandler) Update(c *gin.Context) {
	task, ok := handler.getOwnerTask(c)
	if !ok {
		return
	}

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = model.TaskStatus(req.Status)

	// Sau khi cập nhật thông tin task, lưu lại vào cơ sở dữ liệu
	handler.DB.Save(task)
	c.JSON(http.StatusOK, task)
}

func (handler *TaskHandler) Delete(c *gin.Context) {
	task, ok := handler.getOwnerTask(c)
	if !ok {
		return
	}
	handler.DB.Delete(task)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
