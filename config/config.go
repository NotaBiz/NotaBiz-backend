// Package config provides functionality for loading and managing application
// configuration from environment variables. It uses the `godotenv` package
// to load environment variables from a `.env` file and maps them to a
// `ConfigData` structure defined in the `model` package.
package config

import (
	"NotaBiz-backend/model"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Data is a global variable that holds the application's configuration data.
// It is populated during the initialization of the application.
var Data model.ConfigData

// init initializes the configuration by calling the Load function.
// This ensures that the configuration is loaded as soon as the package is imported.
func init() {
    Load()
}

// Load reads the environment variables from the `.env` file and populates
// the `Data` variable with the application's configuration. It calls
// specific functions to load database, application, and logger configurations.
func Load() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    loadDbConfig()
    loadAppConfig()
    loadLoggerConfig()
    loadServiceConfig()
}

// loadDbConfig loads the database configuration from environment variables
// and assigns them to the `DbConfig` field of the `Data` variable.
// It ensures that all required database environment variables are present.
func loadDbConfig() {
    Data.DbConfig.DbHost = os.Getenv("DB_HOST")
    Data.DbConfig.DbPort = os.Getenv("DB_PORT")
    Data.DbConfig.DbUser = os.Getenv("DB_USER")
    Data.DbConfig.DbPassword = os.Getenv("DB_PASSWORD")
    Data.DbConfig.DbName = os.Getenv("DB_NAME")
    
    if Data.DbConfig.DbHost == "" || Data.DbConfig.DbPort == "" || 
       Data.DbConfig.DbUser == "" || Data.DbConfig.DbPassword == "" || 
       Data.DbConfig.DbName == "" {
        log.Fatal("missing database environment variables")
    }
}

// loadAppConfig loads the application configuration from environment variables
// and assigns them to the `AppConfig` field of the `Data` variable.
// It ensures that all required application environment variables are present
// and validates the `APP_PORT` and `JWT_EXPIRATION` values.
func loadAppConfig() {
    Data.AppConfig.Name = os.Getenv("APP_NAME")
    Data.AppConfig.Version = os.Getenv("APP_VERSION")
    port := os.Getenv("APP_PORT")
    jwtExpiration := os.Getenv("JWT_EXPIRATION")
    Data.AppConfig.JwtSecret = os.Getenv("JWT_SECRET")
    Data.AppConfig.Environment = os.Getenv("APP_ENVIRONMENT")
    Data.AppConfig.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
    Data.AppConfig.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")

    if Data.AppConfig.Name == "" || Data.AppConfig.Version == "" ||
       port == "" || jwtExpiration == "" || Data.AppConfig.JwtSecret == "" ||
       Data.AppConfig.Environment == "" || Data.AppConfig.GoogleClientID == "" ||
       Data.AppConfig.GoogleClientSecret == "" {
        log.Fatal("missing application environment variables")
    }
    
    var err error
    Data.AppConfig.Port, err = strconv.Atoi(port)
    if err != nil {
        log.Fatal("invalid port:", err)
    }
    
    Data.AppConfig.JwtExpiration, err = strconv.Atoi(jwtExpiration)
    if err != nil {
        log.Fatal("invalid jwt expiration:", err)
    }
}

// loadLoggerConfig loads the logger configuration from environment variables
// and assigns them to the `LoggerConfig` field of the `Data` variable.
// It ensures that all required logger environment variables are present
// and validates the numeric values for `LOG_MAX_SIZE`, `LOG_MAX_BACKUPS`,
// and `LOG_MAX_AGE`.
func loadLoggerConfig() {
    Data.LoggerConfig.Path = os.Getenv("LOG_PATH")
    
    var err error
    Data.LoggerConfig.MaxSize, err = strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
    if err != nil {
        log.Fatal("invalid log max size:", err)
    }
    
    Data.LoggerConfig.MaxBackups, err = strconv.Atoi(os.Getenv("LOG_MAX_BACKUPS"))
    if err != nil {
        log.Fatal("invalid log max backups:", err)
    }
    
    Data.LoggerConfig.MaxAge, err = strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
    if err != nil {
        log.Fatal("invalid log max age:", err)
    }
    
    Data.LoggerConfig.Compress = os.Getenv("LOG_COMPRESS") == "true"
}

func loadServiceConfig() {
    Data.ServiceConfig.RedisUrl = os.Getenv("REDIS_URL")
    Data.ServiceConfig.RedisPassword = os.Getenv("REDIS_PASSWORD")
    Data.ServiceConfig.RabbitmqUrl = os.Getenv("RABBITMQ_URL")
    Data.ServiceConfig.SenderEmail = os.Getenv("SENDER_EMAIL")
    Data.ServiceConfig.SenderPassword = os.Getenv("SENDER_PASSWORD")
    Data.ServiceConfig.SmtpHost = os.Getenv("SMTP_HOST")
    smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
    if err != nil {
        log.Fatal("invalid smtp port:", err)
    }
    Data.ServiceConfig.SmtpPort = smtpPort

    if Data.ServiceConfig.RedisUrl == "" || Data.ServiceConfig.RabbitmqUrl == "" || Data.ServiceConfig.SenderEmail == "" || Data.ServiceConfig.SenderPassword == "" || Data.ServiceConfig.SmtpHost == "" || Data.ServiceConfig.SmtpPort == 0 {
        log.Fatal("missing service environment variables")
    }
}