package main

import (
	_ "github.com/dazraf/go-api-example/docs"
	"github.com/dazraf/go-api-example/handlers"
	"github.com/dazraf/go-api-example/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           User API
// @version         1.0
// @description     A simple user management API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// Initialize the user store with some sample data
	userStore := store.NewMemoryUserStore()

	// Add some initial users
	userStore.Create(store.User{Name: "John Doe", Email: "john@example.com"})
	userStore.Create(store.User{Name: "Jane Smith", Email: "jane@example.com"})

	// Create handler with dependency injection
	userHandler := handlers.NewUserHandler(userStore)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", userHandler.GetUsers)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.POST("/users", userHandler.CreateUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	} // Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8080")
}
