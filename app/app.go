package app

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/database"
	"NotaBiz-backend/model"
	"NotaBiz-backend/router"
	"flag"
	"time"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	var configData *model.ConfigData = config.GetConfig()
	// Connect to database
	db, err := database.ConnectDB(configData)
	if err != nil {
		log.Fatal(err)
	}
	// optional flag to seeding data
	seedCommand := flag.Bool("seed", false, "seed the database")
	flag.Parse()
	if *seedCommand {
		database.SeedData(db)
	}
	conn, _ := db.DB()
	defer func() {
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// setup gin
	if configData.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%d", configData.AppConfig.Port)},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           600 * time.Second,
	}))
	logger, err := SetupLogger(configData)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	// logger to catch request response in api
	r.Use(ResponseLogger(logger))

	// setup router
	router.InitRouter(r, string(configData.AppConfig.Version[0]))

	fmt.Printf("Running %s version %s on port %d\n", configData.AppConfig.Name, configData.AppConfig.Version, configData.AppConfig.Port)
	err = r.Run(":" + strconv.Itoa(configData.AppConfig.Port))
	if err != nil {
		log.Fatal(err)
	}
}
