package interfaces

import "github.com/ratheeshkumar25/task-mgt/internal/models"

type TaskRepoInter interface {
	//user repo
	CreateUser(user *models.Users) error
	FindUserByID(userID uint) (*models.Users, error)
	GetUserByUsername(usename string) (*models.Users, error)
	GetUserList() ([]*models.Users, error)

	//task repo
	CreateTask(task *models.Task) error
	GetFilteredTasks(status, dueDateAfter, sortBy, sortOrder string, page, limit int) ([]models.Task, int64, error)
	//GetAllTask() ([]models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id uint) error
}
