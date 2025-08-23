package app

import (
	"NotaBiz-backend/model"
	"NotaBiz-backend/router"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func RunServer() {
	configData, err := initEnv()
	if err != nil {
		log.Fatal(err)
	}
	
	routes := router.InitRouter()

	fmt.Println("Running server on port " + strconv.Itoa(configData.AppConfig.Port))
	err = http.ListenAndServe(":"+strconv.Itoa(configData.AppConfig.Port), routes)
	if err != nil {
		log.Fatal(err)
	}
}

func initEnv() (*model.ConfigData, error) {
	var configData model.ConfigData
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	configData.DbConfig.DbHost = os.Getenv("DB_HOST")
	configData.DbConfig.DbPort = os.Getenv("DB_PORT")
	configData.DbConfig.DbUser = os.Getenv("DB_USER")
	configData.DbConfig.DbPassword = os.Getenv("DB_PASSWORD")
	configData.DbConfig.DbName = os.Getenv("DB_NAME")
	if configData.DbConfig.DbHost == "" || configData.DbConfig.DbPort == "" || configData.DbConfig.DbUser == "" || configData.DbConfig.DbPassword == "" || configData.DbConfig.DbName == "" {
		return nil, errors.New("missing database environment variables")
	}

	port := os.Getenv("APP_PORT")
	jwtExpiration := os.Getenv("JWT_EXPIRATION")
	configData.AppConfig.JwtSecret = os.Getenv("JWT_SECRET")
	if port == "" || jwtExpiration == "" || configData.AppConfig.JwtSecret == "" {
		return nil, errors.New("missing application environment variables")
	}
	configData.AppConfig.Port, err = strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	configData.AppConfig.JwtExpiration, err = strconv.Atoi(jwtExpiration)
	if err != nil {
		return nil, err
	}
	return &configData, nil
}