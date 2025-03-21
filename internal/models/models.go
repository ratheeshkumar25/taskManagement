package models

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "Pending"
	TaskStatusInProgress TaskStatus = "In Progress"
	TaskStatusCompleted  TaskStatus = "Completed"
)

type Users struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
}

type Task struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	UserID      uint       `json:"-"`
}

type TaskWithTotal struct {
	Tasks []Task `json:"tasks"`
	Total int64  `json:"total"`
}

type TaskResponse struct {
	ID        uint   `json:"ID"`
	Title     string `json:"Title"`
	Status    string `json:"Status"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

// IsValidStatus checks if the task status is valid
func IsValidStatus(status TaskStatus) bool {
	return status == TaskStatusPending ||
		status == TaskStatusInProgress ||
		status == TaskStatusCompleted
}
