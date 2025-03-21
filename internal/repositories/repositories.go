package repositories

import (
	"fmt"

	"github.com/ratheeshkumar25/task-mgt/internal/models"
	inter "github.com/ratheeshkumar25/task-mgt/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

// GetUserByUsername services
func (t *TaskRepository) GetUserByUsername(usename string) (*models.Users, error) {
	var user models.Users

	if err := t.DB.Model(&models.Users{}).Where("username = ?", &usename).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserList implements
func (t *TaskRepository) GetUserList() ([]*models.Users, error) {
	var user []*models.Users
	if err := t.DB.Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Create User implements
func (t *TaskRepository) CreateUser(user *models.Users) error {
	if err := t.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// Find user implements
func (t *TaskRepository) FindUserByID(userID uint) (*models.Users, error) {
	var user models.Users
	if err := t.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create Task implements
func (t *TaskRepository) CreateTask(task *models.Task) error {
	if err := t.DB.Create(&task).Error; err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) GetFilteredTasks(status, dueDateAfter, sortBy, sortOrder string, page, limit int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	db := r.DB.Model(&models.Task{})

	// Filters
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if dueDateAfter != "" {
		db = db.Where("due_date >= ?", dueDateAfter)
	}

	// Get total count before pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if sortBy != "" && sortOrder != "" {
		db = db.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	}

	// Pagination
	offset := (page - 1) * limit
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// Final query
	if err := db.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// // GetAll Task implements
// func (t *TaskRepository) GetAllTask() ([]models.Task, error) {
// 	var tasks []models.Task
// 	if err := t.DB.Find(&tasks).Error; err != nil {
// 		return nil, err
// 	}
// 	return tasks, nil
// }

// Get Task ByID implements
func (t *TaskRepository) GetTaskByID(id uint) (*models.Task, error) {
	var tasks models.Task
	if err := t.DB.First(&tasks, id).Error; err != nil {
		return nil, err
	}
	return &tasks, nil
}

// Update Task implements
func (t *TaskRepository) UpdateTask(task *models.Task) error {
	if err := t.DB.Save(&task).Error; err != nil {
		return err
	}
	return nil
}

// Delete Task implements
func (r *TaskRepository) DeleteTask(id uint) error {
	result := r.DB.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// NewTaskRepository: Constructor function
func NewTaskRepository(db *gorm.DB) inter.TaskRepoInter {
	return &TaskRepository{
		DB: db,
	}
}
