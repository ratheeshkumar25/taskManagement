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
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockTaskRepoInter(ctrl)
	logger := log.Default()
	service := services.NewTaskService(repoMock, nil, logger)

	user := &models.Users{Username: "testuser", PasswordHash: "password"}
	repoMock.EXPECT().CreateUser(user).Return(nil)

	err := service.CreateUser(user)
	assert.NoError(t, err)
}

func TestCreateUser_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockTaskRepoInter(ctrl)
	logger := log.Default()
	service := services.NewTaskService(repoMock, nil, logger)

	user := &models.Users{Username: "testuser", PasswordHash: "password"}
	repoMock.EXPECT().CreateUser(user).Return(errors.New("db error"))

	err := service.CreateUser(user)
	assert.Error(t, err)
}

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
