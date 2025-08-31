package database

import (
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/entity"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config *model.ConfigData) (*gorm.DB, error) {
	err := createDB(config)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user= %s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.DbConfig.DbHost, config.DbConfig.DbUser, config.DbConfig.DbPassword, strings.ToLower(config.DbConfig.DbName), config.DbConfig.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}

	// migrate db
	err = db.AutoMigrate(
		&entity.MasterRole{},
		&entity.MasterSubscription{},
		&entity.MasterCompany{},
		&entity.MasterUser{},
		&entity.Feedback{},
		&entity.Transaction{},
		&entity.Attachment{},
	)
	if err != nil {
		log.Println("Error migrating database:", err)
		return nil, err
	}

	return db, nil
}

func createDB(config *model.ConfigData) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.DbConfig.DbHost, config.DbConfig.DbUser, config.DbConfig.DbPassword, config.DbConfig.DbPort)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err = db.QueryRow(query, strings.ToLower(config.DbConfig.DbName)).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		log.Println("Database already exists, avoiding create DB")
		return nil
	}

	// Create DB
	_, err = db.Exec("CREATE DATABASE " + config.DbConfig.DbName)
	if err != nil {
		return err
	}
	log.Println("Database created successfully")

	return nil
}

func SeedData(db *gorm.DB) error {
	err := db.Create(&RoleSeed).Error
	if err != nil {
		return err
	}
	log.Println("Seed roles data successfully")
	
	err = db.Create(&SubscriptionSeed).Error
	if err != nil {
		return err
	}
	log.Println("Seed subscription data successfully")

	err = db.Create(&CompanySeed).Error
	if err != nil {
		return err
	}
	log.Println("Seed company data successfully")

	err = db.Create(&UserSeed).Error
	if err != nil {
		return err
	}
	log.Println("Seed user data successfully")

	return nil
}
