package config

import (
	"NotaBiz-backend/database"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis_rate/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	Redis    *redis.Client
	Limiter  *redis_rate.Limiter
	RabbitMQ *amqp.Connection
	Channel  *amqp.Channel
)

func InitServices() error {
	// Connect to the database
	db, err := database.ConnectDB(&Data)
	if err != nil {
		return fmt.Errorf("database connection failed: %v", err)
	}
	DB = db
	sqlDb, _ := db.DB()

	//set max conn, idle and lifetime for all connections in gorm.DB object
	sqlDb.SetConnMaxLifetime(Data.DbConfig.MaxLifeTime)
	sqlDb.SetConnMaxIdleTime(time.Duration(Data.DbConfig.MaxIdle))
	sqlDb.SetMaxOpenConns(Data.DbConfig.MaxConn)
	sqlDb.SetMaxIdleConns(Data.DbConfig.MaxIdle)

	// Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:         Data.ServiceConfig.RedisUrl,
		Password:     Data.ServiceConfig.RedisPassword,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %v", err)
	}
	Redis = redisClient

	rdb := redis.NewClient(&redis.Options{
		Addr: Data.ServiceConfig.RedisUrl,
	})
	_ = rdb.FlushDB(ctx).Err()

	Limiter = redis_rate.NewLimiter(rdb)

	// RabbitMQ connection
	rabbitmqURL := Data.ServiceConfig.RabbitmqUrl
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return fmt.Errorf("rabbitmq connection failed: %v", err)
	}
	RabbitMQ = conn

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel failed: %v", err)
	}
	Channel = ch

	// Declare queues
	if err := setupQueues(ch); err != nil {
		return fmt.Errorf("queue setup failed: %v", err)
	}

	return nil
}

func setupQueues(ch *amqp.Channel) error {
	// Declare exchanges
	err := ch.ExchangeDeclare(
		"otp.exchange", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	// Declare queues
	queues := []string{"otp.email", "otp.whatsapp", "otp.retry"}

	for _, queueName := range queues {
		_, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			amqp.Table{
				"x-message-ttl":          300000, // 5 minutes TTL
				"x-max-retries":          3,
				"x-dead-letter-exchange": "otp.dlx",
			},
		)
		if err != nil {
			return err
		}

		// Bind queue to exchange
		err = ch.QueueBind(
			queueName,      // queue name
			queueName,      // routing key
			"otp.exchange", // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	// Dead letter queue
	_, err = ch.QueueDeclare("otp.dlq", true, false, false, false, nil)
	return err
}

func Cleanup() {
	if conn, _ := DB.DB(); conn != nil {
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}
	if Redis != nil {
		if err := Redis.Close(); err != nil {
			log.Fatal(err)
		}
	}
	if Channel != nil {
		if err := Channel.Close(); err != nil {
			log.Fatal(err)
		}
	}
	if RabbitMQ != nil {
		if err := RabbitMQ.Close(); err != nil {
			log.Fatal(err)
		}
	}
}