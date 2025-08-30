package model

type (
	ConfigData struct {
		DbConfig
		AppConfig
		LoggerConfig
	}

	DbConfig struct {
		DbHost     string
		DbPort     string
		DbUser     string
		DbPassword string
		DbName     string
	}

	AppConfig struct {
		Name          string
		Version       string
		Port          int
		Environment   string
		JwtSecret     string
		JwtExpiration int
	}

	LoggerConfig struct {
		Path       string
		MaxSize    int
		MaxBackups int
		MaxAge     int
		Compress   bool
	}
)
