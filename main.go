// Package main implements the entry point for the NotaBiz-backend application.
// It sets the local time zone to "Asia/Jakarta" and starts the server.
package main

import (
	"NotaBiz-backend/app"
	"NotaBiz-backend/docs"
	"time"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	docs.SwaggerInfo.Title = "NotaBiz API"
	docs.SwaggerInfo.Description = "API documentation for NotaBiz application."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}
	app.RunServer()
}