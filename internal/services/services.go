package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ratheeshkumar25/task-mgt/config"
	"github.com/ratheeshkumar25/task-mgt/internal/models"
	repoIface "github.com/ratheeshkumar25/task-mgt/internal/repositories/interfaces"
	inter "github.com/ratheeshkumar25/task-mgt/internal/services/interfaces"
	"github.com/ratheeshkumar25/task-mgt/utility"
	"golang.org/x/crypto/bcrypt"
)

type TaskServices struct {
	Repo   repoIface.TaskRepoInter
	redis  *config.RedisService
	Logger *log.Logger
}

// LoginUser: Authenticates user and generates JWT token
func (t *TaskServices) LoginUser(username string, password string) (string, error) {
	foundUser, err := t.Repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	cfg := config.LoadConfig()
	token, err := utility.GenerateToken(cfg.SECERETKEY, foundUser.Username, foundUser.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateUser: Registers a new user
func (t *TaskServices) CreateUser(user *models.Users) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return t.Repo.CreateUser(user)
}

// CreateTask: Creates a task and clears Redis cache asynchronously
func (t *TaskServices) CreateTask(task *models.Task) error {
	err := t.Repo.CreateTask(task)
	if err == nil {
		go func() {
			if delErr := t.redis.DeleteFromRedis("tasks_list"); delErr != nil {
				t.Logger.Println("Redis delete error:", delErr)
			}
		}()
	}
	return err
}

// GetAllTasks: Fetches tasks with caching support
func (t *TaskServices) GetAllTasks(status string, dueDateAfter string, sortBy string, sortOrder string, page int, limit int) ([]models.Task, int64, error) {
	cacheKey := fmt.Sprintf("tasks_list:status=%s:dueAfter=%s:sortBy=%s:order=%s:page=%d:limit=%d",
		status, dueDateAfter, sortBy, sortOrder, page, limit)

	// Try to get from Redis
	cachedData, err := t.redis.GetFromRedis(cacheKey)
	if err == nil && cachedData != "" {
		var cachedTasks []models.TaskWithTotal
		if jsonErr := json.Unmarshal([]byte(cachedData), &cachedTasks); jsonErr == nil && len(cachedTasks) > 0 {
			return cachedTasks[0].Tasks, cachedTasks[0].Total, nil
		}
	}

	// Fetch from DB if not cached
	tasks, total, err := t.Repo.GetFilteredTasks(status, dueDateAfter, sortBy, sortOrder, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Store result in Redis
	go func() {
		cacheData := []models.TaskWithTotal{{Tasks: tasks, Total: total}}
		jsonData, _ := json.Marshal(cacheData)
		_ = t.redis.SetDataInRedis(cacheKey, jsonData, 5*time.Minute)
	}()

	return tasks, total, nil
}

// GetTaskByID: Retrieves a single task, caches it for quick access
func (t *TaskServices) GetTaskByID(id uint) (*models.Task, error) {
	key := fmt.Sprintf("task_%d", id)

	// Try to fetch from cache
	cachedData, err := t.redis.GetFromRedis(key)
	if err == nil && cachedData != "" {
		var task models.Task
		json.Unmarshal([]byte(cachedData), &task)
		return &task, nil
	}

	// Fetch from DB if not found in cache
	task, err := t.Repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the fetched task
	go func() {
		jsonData, _ := json.Marshal(task)
		_ = t.redis.SetDataInRedis(key, jsonData, 5*time.Minute)
	}()

	return task, nil
}

// UpdateTask: Updates a task and clears Redis cache concurrently
func (t *TaskServices) UpdateTask(task *models.Task) error {
	err := t.Repo.UpdateTask(task)
	if err == nil {
		go func(id uint) {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				if delErr := t.redis.DeleteFromRedis("tasks_list"); delErr != nil {
					t.Logger.Println("Redis delete error:", delErr)
				}
			}()

			go func() {
				defer wg.Done()
				if delErr := t.redis.DeleteFromRedis(fmt.Sprintf("task_%d", id)); delErr != nil {
					t.Logger.Println("Redis delete error:", delErr)
				}
			}()

			wg.Wait()
		}(task.ID)
	}
	return err
}

// DeleteTask: Deletes a task and clears Redis cache concurrently
func (t *TaskServices) DeleteTask(id uint) error {
	err := t.Repo.DeleteTask(id)
	if err == nil {
		go func(id uint) {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				if delErr := t.redis.DeleteFromRedis("tasks_list"); delErr != nil {
					t.Logger.Println("Redis delete error:", delErr)
				}
			}()

			go func() {
				defer wg.Done()
				if delErr := t.redis.DeleteFromRedis(fmt.Sprintf("task_%d", id)); delErr != nil {
					t.Logger.Println("Redis delete error:", delErr)
				}
			}()

			wg.Wait()
		}(id)
	}
	return err
}

// NewTaskService: Constructor function
func NewTaskService(repo repoIface.TaskRepoInter, redis *config.RedisService, logger *log.Logger) inter.TaskServiceInter {
	return &TaskServices{
		Repo:   repo,
		redis:  redis,
		Logger: logger,
	}
}
