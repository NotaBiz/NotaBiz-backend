// Package app provides the functionality to initialize, configure, and run
// the Gin web server for the NotaBiz-backend application.
//
// It handles server setup, middleware registration, CORS configuration,
// logging, database connection management, and routing initialization.
// The package is designed to be the entry point for running the backend
// service.
package app

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/middleware"
	"NotaBiz-backend/model"
	"NotaBiz-backend/router"
	"NotaBiz-backend/utils"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RunServer initializes and starts the Gin web server.
//
// RunServer performs the following tasks:
//   - Loads application configuration from the config package.
//   - Establishes a connection to the database and optionally seeds data
//     if the `-seed` flag is provided at runtime.
//   - Configures Gin to run in release mode if the environment is set
//     to production.
//   - Registers CORS, security headers, and custom logging middleware.
//   - Initializes HTTP routes via the router package.
//   - Starts the HTTP server on the configured port.
//
// The function blocks the main goroutine until the server shuts down or
// encounters a fatal error.
func RunServer() {
	// Load configuration
	var configData *model.ConfigData = &config.Data

	// init services
	err := config.InitServices()
	defer config.Cleanup()
	if err != nil {
		log.Fatal("Failed to initialize services:", err)
	}

	// Start background workers
	go utils.StartOTPWorkers(3) // 3 concurrent workers, each for email and whatsapp
	go utils.RunScheduler()

	// Configure Gin mode based on environment
	if configData.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin engine
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%d", configData.AppConfig.Port)},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           600 * time.Second,
	}))

	// Rate size limiter
	r.Use(middleware.MaxSizeMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	logger, err := middleware.SetupLogger(configData)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	r.Use(middleware.ResponseLogger(logger))
	r.Use(middleware.CORSMiddleware())

	// Initialize routes
	router.InitRouter(r)

	// Start server
	fmt.Printf("Running %s version %s on port %d\n",
		configData.AppConfig.Name,
		configData.AppConfig.Version,
		configData.AppConfig.Port)

	err = r.Run(":" + strconv.Itoa(configData.AppConfig.Port))
	if err != nil {
		log.Fatal(err)
	}
}