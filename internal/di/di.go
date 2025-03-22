package di

import (
	"github.com/gin-gonic/gin"
	"github.com/ratheeshkumar25/task-mgt/config"
	"github.com/ratheeshkumar25/task-mgt/internal/db"
	"github.com/ratheeshkumar25/task-mgt/internal/handlers"
	"github.com/ratheeshkumar25/task-mgt/internal/repositories"
	"github.com/ratheeshkumar25/task-mgt/internal/services"
	"github.com/ratheeshkumar25/task-mgt/utility"
)

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

	v1 := router.Group("/api/v1")

	// Inject Handler Layer
	handlers.NewTaskHandler(v1, taskService, cfg.SECERETKEY, redisClient.Client)

	// Start server
	if err := router.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
