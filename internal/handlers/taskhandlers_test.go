package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/ratheeshkumar25/task-mgt/internal/handlers"
	"github.com/ratheeshkumar25/task-mgt/internal/mocks"
	"github.com/ratheeshkumar25/task-mgt/internal/models"
	"github.com/stretchr/testify/assert"
)

// Setup function for initializing gin context and router
func setupTestRouter() (*gin.Engine, *gin.RouterGroup) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	apiGroup := router.Group("/api/v1")
	return router, apiGroup
}

// Test Register Handler
func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInter(ctrl)
	router, apiGroup := setupTestRouter()

	h := handlers.TaskHandler{SVC: mockService, JWTSecret: "secret"}
	apiGroup.POST("/register", h.Register)

	user := models.Users{
		Username:     "test@example.com",
		PasswordHash: "hashedpassword",
	}

	mockService.EXPECT().CreateUser(gomock.Any()).Return(nil)

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "user registered successfully")
}

// Test Login Handler
func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInter(ctrl)
	router, apiGroup := setupTestRouter()

	h := handlers.TaskHandler{SVC: mockService, JWTSecret: "secret"}
	apiGroup.POST("/login", h.Login)

	loginData := map[string]string{"username": "test@example.com", "password": "password123"}
	mockService.EXPECT().LoginUser("test@example.com", "password123").Return("mockToken", nil)

	reqBody, _ := json.Marshal(loginData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mockToken")
}

// Test Create Task
func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInter(ctrl)
	router, apiGroup := setupTestRouter()

	h := handlers.TaskHandler{SVC: mockService, JWTSecret: "secret"}
	apiGroup.POST("/tasks", h.CreateTask)

	task := models.Task{
		Title:     "Test Task",
		Status:    models.TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.EXPECT().CreateTask(gomock.Any()).Return(nil)

	reqBody, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test Task")
}

// // Test Get All Tasks
func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInter(ctrl)
	router, apiGroup := setupTestRouter()

	h := handlers.TaskHandler{SVC: mockService, JWTSecret: "secret"}
	apiGroup.GET("/tasks", h.GetAllTasks)

	mockTasks := []models.Task{
		{ID: 1, Title: "Task 1", Status: models.TaskStatusCompleted, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Task 2", Status: models.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	mockService.EXPECT().GetAllTasks("", "", "due_date", "asc", 1, 10).Return(mockTasks, int64(2), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task 1")
	assert.Contains(t, w.Body.String(), "Task 2")
}

// Test Get Task by ID
func TestGetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInter(ctrl)
	router, apiGroup := setupTestRouter()

	h := handlers.TaskHandler{SVC: mockService, JWTSecret: "secret"}
	apiGroup.GET("/tasks/:id", h.GetTaskByID)

	mockTask := models.Task{ID: 1, Title: "Task 1", Status: models.TaskStatusPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mockService.EXPECT().GetTaskByID(uint(1)).Return(&mockTask, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task 1")
}

// Test Update Task
func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInter(ctrl)
	router, apiGroup := setupTestRouter()

	h := handlers.TaskHandler{SVC: mockService, JWTSecret: "secret"}
	apiGroup.PUT("/tasks/:id", h.UpdateTask)

	updatedTask := models.Task{ID: 1, Title: "Updated Task", Status: models.TaskStatusCompleted, UpdatedAt: time.Now()}

	mockService.EXPECT().UpdateTask(gomock.Any()).Return(nil)

	reqBody, _ := json.Marshal(updatedTask)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Task")
}

// Test Delete Task
func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockTaskServiceInter(ctrl)
	router, apiGroup := setupTestRouter()

	h := handlers.TaskHandler{SVC: mockService, JWTSecret: "secret"}
	apiGroup.DELETE("/tasks/:id", h.DeleteTask)

	mockService.EXPECT().GetTaskByID(uint(1)).Return(&models.Task{ID: 1}, nil)
	mockService.EXPECT().DeleteTask(uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "task deleted successfully")
}
