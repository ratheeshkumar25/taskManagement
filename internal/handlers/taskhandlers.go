package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/ratheeshkumar25/task-mgt/internal/middleware"
	"github.com/ratheeshkumar25/task-mgt/internal/models"
	inter "github.com/ratheeshkumar25/task-mgt/internal/services/interfaces"
)

type TaskHandler struct {
	SVC       inter.TaskServiceInter
	JWTSecret string
}

func NewTaskHandler(router *gin.RouterGroup, svc inter.TaskServiceInter, secret string, redisClient *redis.Client) {
	h := &TaskHandler{SVC: svc, JWTSecret: secret}

	// Public routes
	router.POST("/login", h.Login)
	router.POST("/register", h.Register)

	// Protected task routes
	auth := router.Group("/tasks")
	auth.Use(
		middleware.AuthMiddleware(secret),
		middleware.RateLimitMiddleware(redisClient, 60, time.Minute),
	)
	{
		auth.POST("", h.CreateTask)
		auth.GET("", h.GetAllTasks)
		auth.GET("/:id", h.GetTaskByID)
		auth.PUT("/:id", h.UpdateTask)
		auth.DELETE("/:id", h.DeleteTask)
	}
}

func (h *TaskHandler) Register(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err := h.SVC.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *TaskHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	token, err := h.SVC.LoginUser(loginData.Username, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil || task.Title == "" || !models.IsValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input or status"})
		return
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if err := h.SVC.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	// Read query params
	status := c.Query("status")
	dueDateAfter := c.Query("due_date_after")
	sortBy := c.DefaultQuery("sort_by", "due_date")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Call service
	tasks, total, err := h.SVC.GetAllTasks(status, dueDateAfter, sortBy, sortOrder, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Map tasks to response DTO
	var taskResponses []models.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, models.TaskResponse{
			ID:        task.ID,
			Title:     task.Title,
			Status:    string(task.Status),
			CreatedAt: task.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: task.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	// Final response
	c.JSON(http.StatusOK, gin.H{
		"tasks": taskResponses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	task, err := h.SVC.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil || task.Title == "" || !models.IsValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	task.ID = uint(id)
	task.UpdatedAt = time.Now()
	if err := h.SVC.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	// Check if task exists first
	task, err := h.SVC.GetTaskByID(uint(id))
	if err != nil || task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// Proceed to delete
	if err := h.SVC.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
