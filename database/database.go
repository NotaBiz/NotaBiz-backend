// Package database provides functionality for connecting to a PostgreSQL database,
// creating the database if it does not exist, running migrations, and seeding initial data.
package database

import (
	"NotaBiz-backend/model"
	"NotaBiz-backend/model/entity"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ConnectDB establishes a connection to the PostgreSQL database using the provided
// configuration. It creates the database if it does not exist, runs migrations for
// the defined entities, and returns a GORM database instance.
//
// Parameters:
//   - config: A pointer to the ConfigData structure containing database configuration.
//
// Returns:
//   - A pointer to the GORM database instance.
//   - An error if the connection or migration fails.
func ConnectDB(config *model.ConfigData) (*gorm.DB, error) {
	flushCommand := flag.Bool("flush", false, "flush the database")
	seedCommand := flag.Bool("seed", false, "seed the database")
	flag.Parse()
	if *flushCommand {
		err := flushDB(config)
		if err != nil {
			return nil, err
		}
		*seedCommand = true
	} else {
		err := createDB(config)
		if err != nil {
			return nil, err
		}
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

	if *seedCommand {
		err := seedData(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// createDB checks if the database specified in the configuration exists.
// If it does not exist, it creates the database.
//
// Parameters:
//   - config: A pointer to the ConfigData structure containing database configuration.
//
// Returns:
//   - An error if the database creation fails.
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

	_, err = db.Exec("CREATE DATABASE " + config.DbConfig.DbName)
	if err != nil {
		return err
	}
	log.Println("Database created successfully")

	return nil
}

// SeedData seeds the database with initial data for roles, subscriptions, companies, and users.
//
// Parameters:
//   - db: A pointer to the GORM database instance.
//
// Returns:
//   - An error if any of the seeding operations fail.
func seedData(db *gorm.DB) error {
	err := db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&RoleSeed).Error
	if err != nil {
		return err
	}
	log.Println("Seed roles data successfully")

	err = db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&SubscriptionSeed).Error
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

func flushDB(config *model.ConfigData) error {
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
		_, err = db.Exec("DROP DATABASE " + config.DbConfig.DbName)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec("CREATE DATABASE " + config.DbConfig.DbName)
	if err != nil {
		return err
	}
	log.Println("Database created successfully")

	return nil
}
