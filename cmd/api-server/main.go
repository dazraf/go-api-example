package main

import (
	"log"

	"github.com/dazraf/go-api-example/internal/app"
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
	// Initialize application
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start server
	log.Printf("Starting server on %s", application.Config.Server.Address)
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
