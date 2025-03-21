package interfaces

import "github.com/ratheeshkumar25/task-mgt/internal/models"

type TaskServiceInter interface {
	//Service to handle the user
	CreateUser(user *models.Users) error
	LoginUser(username string, password string) (string, error)
	//Service to handle the tasks
	CreateTask(task *models.Task) error
	GetAllTasks(status, dueDateAfter, sortBy, sortOrder string, page, limit int) ([]models.Task, int64, error)
	GetTaskByID(id uint) (*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id uint) error
}
