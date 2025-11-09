package utils

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/model"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartOTPWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		// auto restart worker
		go func(workerID int) {
			for {
				log.Printf("Starting worker email %d", workerID)
				otpWorker("otp.email")
				log.Printf("Worker email %d died, restarting in 5 seconds", workerID)
				time.Sleep(5 * time.Second)
			}
		}(i)
		go func(workerID int) {
			for {
				log.Printf("Starting worker whatsapp %d", workerID)
				otpWorker("otp.whatsapp")
				log.Printf("Worker whatsapp %d died, restarting in 5 seconds", workerID)
				time.Sleep(5 * time.Second)
			}
		}(i)
	}
}

func otpWorker(queueName string) {
	msgs, err := config.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Failed to register consumer for %s: %v", queueName, err)
		return
	}

	// Infinite loop to consume every messages send, loop continues until channel is closed
	// If queue empty, worker will sleep/wait for a message
	// When message arrives, worker will process it, acknowledge it, then wait for next message
	for msg := range msgs {
		var otpMsg model.OTPMessage
		if err := json.Unmarshal(msg.Body, &otpMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			msg.Nack(false, false)
			continue
		}

		// Process message with retry logic
		success := ProcessOTPMessage(otpMsg)

		if success {
			msg.Ack(false)
			log.Printf("Successfully sent OTP via %s to %s%s",
				otpMsg.Method, otpMsg.Email, otpMsg.PhoneNumber)
		} else {
			// Retry logic
			if otpMsg.Retries < 3 {
				otpMsg.Retries++
				time.Sleep(time.Duration(otpMsg.Retries) * time.Second) // Exponential backoff
				PublishMessage("otp.retry", otpMsg)
				msg.Ack(false)
			} else {
				log.Printf("Failed to send OTP after 3 retries: %+v", otpMsg)
				msg.Nack(false, false) // Send to DLQ
			}
		}
	}
}

func ProcessOTPMessage(msg model.OTPMessage) bool {
	switch msg.Method {
	case "email":
		return sendOTPByEmail(msg.Email, msg.OTP) == nil
	case "whatsapp":
		return sendOTPByWhatsApp(msg.PhoneNumber, msg.OTP) == nil
	default:
		return false
	}
}

func PublishMessage(queueName string, message model.OTPMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return config.Channel.Publish(
		"otp.exchange", // exchange
		queueName,      // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		})
}

func GetOTPKey(email, phoneNumber, method string) string {
	if method == "email" {
		return fmt.Sprintf("otp:email:%s", email)
	}
	return fmt.Sprintf("otp:whatsapp:%s", phoneNumber)
}

func GenerateOTP() string {
	max := big.NewInt(999999)
	min := big.NewInt(100000)
	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return strconv.Itoa(int(time.Now().Unix()%900000 + 100000))
	}
	return strconv.Itoa(int(n.Int64() + min.Int64()))
}

func sendOTPByEmail(email, otp string) error {
	smtpHost := config.Data.ServiceConfig.SmtpHost
	smtpPort := config.Data.ServiceConfig.SmtpPort
	senderEmail := config.Data.ServiceConfig.SenderEmail
	senderPassword := config.Data.ServiceConfig.SenderPassword

	subject := "NotaBiz Registration Verification"
	body := fmt.Sprintf("Your OTP is: %s\n\nThis OTP will expire in 5 minutes.", otp)
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", email, subject, body)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, senderEmail, []string{email}, []byte(message))

	if err != nil {
		log.Printf("Email send error: %v", err)
		return err
	}

	return nil
}

func sendOTPByWhatsApp(phoneNumber, otp string) error {
	message := fmt.Sprintf("Notabiz Registration Verification. \nYour OTP is: %s. This OTP will expire in 5 minutes.", otp)
	log.Printf("Sending WhatsApp to %s: %s", phoneNumber, message)

	request := map[string]string{
		"target":  phoneNumber,
		"message": message,
	}
	reqJson, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Failed to marshal WhatsApp request: %v \n", err)
		return err
	}

	req, err := http.NewRequest("POST", "https://api.fonnte.com/send", bytes.NewBuffer(reqJson))
	if err != nil {
		fmt.Printf("Failed to send WhatsApp request: %v \n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.Data.ServiceConfig.FonnteToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var resBody model.WhatsAppResponse
		if err := json.NewDecoder(res.Body).Decode(&resBody); err == nil {
			fmt.Printf("WhatsApp API error: %t - %s \n", resBody.Status, resBody.Reason)
		}
		fmt.Printf("WhatsApp API returned non-200 status: %s \n", resBody.Reason)
		return fmt.Errorf("WhatsApp API error: %s", resBody.Reason)
	}

	return nil
}
