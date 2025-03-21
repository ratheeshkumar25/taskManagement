package di

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ratheeshkumar25/task-mgt/cmd/docs"
	"github.com/ratheeshkumar25/task-mgt/config"

	"github.com/ratheeshkumar25/task-mgt/internal/db"
	"github.com/ratheeshkumar25/task-mgt/internal/handlers"
	"github.com/ratheeshkumar25/task-mgt/internal/repositories"
	"github.com/ratheeshkumar25/task-mgt/internal/services"
	"github.com/ratheeshkumar25/task-mgt/utility"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Task Management System
// @version 1.0
// @description A simple task management system built with Golang.
// @termsOfService http://swagger.io/terms/

// @contact.name Ratheesh Kumar
// @license.name Apache 2.0

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func Init() {

	log := utility.InitLogger()
	log.Println("Starting Task Management API Server...🔥")
	// Load config
	cfg := config.LoadConfig()

	// Initialize Redis
	redisClient, err := config.SetupRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// Initialize Database
	dbConn := db.ConnectDB(cfg)

	// Initialize Repository
	taskRepo := repositories.NewTaskRepository(dbConn)

	// Initialize Service Layer
	taskService := services.NewTaskService(taskRepo, redisClient, log)

	// Initialize Router
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")

	// Inject Handler Layer — pass only redisClient.Client
	handlers.NewTaskHandler(v1, taskService, cfg.SECERETKEY, redisClient.Client)

	// Start server
	if err := router.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

}
