package services_test

import (
	"errors"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ratheeshkumar25/task-mgt/internal/mocks"
	"github.com/ratheeshkumar25/task-mgt/internal/models"
	"github.com/ratheeshkumar25/task-mgt/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Create user test case
func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockTaskRepoInter(ctrl)
	logger := log.Default()
	service := services.NewTaskService(repoMock, nil, logger)

	user := &models.Users{Username: "testuser@example.com", PasswordHash: "password"}

	// Mock `GetUserByUsername` to return nil (user not found)
	repoMock.EXPECT().GetUserByUsername(user.Username).Return(nil, gorm.ErrRecordNotFound)

	// Mock `CreateUser`
	repoMock.EXPECT().CreateUser(gomock.Any()).Return(nil)

	err := service.CreateUser(user)
	assert.NoError(t, err)
}

func TestCreateUser_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockTaskRepoInter(ctrl)
	logger := log.Default()
	service := services.NewTaskService(repoMock, nil, logger)

	user := &models.Users{Username: "testuser@example.com", PasswordHash: "password"}

	// First mock GetUserByUsername - return "not found" error
	repoMock.EXPECT().
		GetUserByUsername(user.Username).
		Return(nil, errors.New("user not found")).
		Times(1)

	// Then mock CreateUser - return database error
	repoMock.EXPECT().
		CreateUser(gomock.Any()).
		Return(errors.New("db error")).
		Times(1)

	// Run the test
	err := service.CreateUser(user)

	// Assert the expected error occurs
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

// LOgin  user test case
func TestLoginUser_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockTaskRepoInter(ctrl)
	logger := log.Default()
	service := services.NewTaskService(repoMock, nil, logger)

	user := &models.Users{Username: "testuser", PasswordHash: "$2a$10$7s6u.zw6v7v6B47Gf3.YSO.V3eP3G/ZB7yX/WeYBoR2eyEDyc0J5e"}

	repoMock.EXPECT().GetUserByUsername("testuser").Return(user, nil)

	token, err := service.LoginUser("testuser", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)
}
