package model

type (
	ConfigData struct {
		DbConfig
		AppConfig
	}

	DbConfig struct {
		DbHost     string
		DbPort     string
		DbUser     string
		DbPassword string
		DbName     string
	}

	AppConfig struct {
		Port int
		JwtSecret string
		JwtExpiration int
	}
)