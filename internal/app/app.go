package app

import (
	"github.com/dazraf/go-api-example/internal/config"
	"github.com/dazraf/go-api-example/internal/handlers"
	"github.com/dazraf/go-api-example/internal/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/dazraf/go-api-example/api" // Load swagger docs
) // Application holds the application dependencies and configuration
type Application struct {
	Config      *config.Config
	Router      *gin.Engine
	UserStore   store.UserStore
	UserHandler *handlers.UserHandler
}

// New creates and initializes a new application instance
func New() (*Application, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize the user store
	userStore := store.NewMemoryUserStore()

	// Add some initial sample data
	_, _ = userStore.Create(store.User{Name: "John Doe", Email: "john@example.com"})
	_, _ = userStore.Create(store.User{Name: "Jane Smith", Email: "jane@example.com"})

	// Create handler with dependency injection
	userHandler := handlers.NewUserHandler(userStore)

	// Setup router
	router := setupRouter(userHandler, cfg)

	return &Application{
		Config:      cfg,
		Router:      router,
		UserStore:   userStore,
		UserHandler: userHandler,
	}, nil
}

// Run starts the application server
func (a *Application) Run() error {
	return a.Router.Run(a.Config.Server.Address)
}

// setupRouter configures the gin router with all routes and middleware
func setupRouter(userHandler *handlers.UserHandler, cfg *config.Config) *gin.Engine {
	// Set gin mode based on config
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", userHandler.GetUsers)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.POST("/users", userHandler.CreateUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// Swagger endpoint (only in non-production)
	if cfg.Environment != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
