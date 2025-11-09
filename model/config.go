// Package model defines the structures used in the application.
package model

import "time"

type (
	ConfigData struct {
		DbConfig
		AppConfig
		LoggerConfig
		ServiceConfig
	}

	DbConfig struct {
		DbHost      string
		DbPort      string
		DbUser      string
		DbPassword  string
		DbName      string
		MaxIdle     int
		MaxIdleTime time.Duration
		MaxConn     int
		MaxLifeTime time.Duration
	}

	AppConfig struct {
		Name               string
		Version            string
		Port               int
		Environment        string
		JwtSecret          string
		JwtExpiration      int
		GoogleClientID     string
		GoogleClientSecret string
	}

	LoggerConfig struct {
		Path       string
		MaxSize    int
		MaxBackups int
		MaxAge     int
		Compress   bool
	}

	ServiceConfig struct {
		RedisUrl       string
		RedisPassword  string
		RabbitmqUrl    string
		SmtpHost       string
		SmtpPort       int
		SenderEmail    string
		SenderPassword string
		FonnteToken    string
	}
)
