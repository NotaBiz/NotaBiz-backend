package config

import (
	"NotaBiz-backend/model"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Data model.ConfigData

func init() {
	Load()
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	loadDbConfig()
	loadAppConfig()
	loadLoggerConfig()
}

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

func loadAppConfig() {
	Data.AppConfig.Name = os.Getenv("APP_NAME")
	Data.AppConfig.Version = os.Getenv("APP_VERSION")
	port := os.Getenv("APP_PORT")
	jwtExpiration := os.Getenv("JWT_EXPIRATION")
	Data.AppConfig.JwtSecret = os.Getenv("JWT_SECRET")
	Data.AppConfig.Environment = os.Getenv("APP_ENVIRONMENT")
	
	if Data.AppConfig.Name == "" || Data.AppConfig.Version == "" || 
	   port == "" || jwtExpiration == "" || Data.AppConfig.JwtSecret == "" || 
	   Data.AppConfig.Environment == "" {
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

// Getter functions for better encapsulation (optional)
func GetConfig() *model.ConfigData {
	return &Data
}