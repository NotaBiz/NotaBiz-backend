package app

import (
	"NotaBiz-backend/model"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func SetupLogger(configData *model.ConfigData) (logging *zap.Logger, err error) {
	logger := lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/log_%s.log", configData.LoggerConfig.Path, time.Now().Format("2006-01-02")),
		MaxSize:    configData.LoggerConfig.MaxSize,
		MaxBackups: configData.LoggerConfig.MaxBackups,
		MaxAge:     configData.LoggerConfig.MaxAge,
		Compress:   configData.LoggerConfig.Compress,
	}
	var level zapcore.Level
	var encoderConfig zapcore.EncoderConfig
	if configData.AppConfig.Environment == "development" {
		level = zapcore.DebugLevel
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		level = zapcore.InfoLevel
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(&logger), level)

	// Create and return the logger
	logging = zap.New(core, zap.AddCaller())
	return logging, nil
}

func ResponseLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the original request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Failed to read request body", zap.Error(err))
			c.Next()
			return
		}

		// Restore the request body so it can be read again
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Create a custom response writer to capture response data
		writer := &responseWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = writer

		// Process the request
		c.Next()

		// Extract request body as a string
		requestBody := string(bodyBytes)
		user, _ := c.Get("user")
		userData, ok := user.(*model.JwtClaims)
		if !ok {
			if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" && writer.statusCode != 200 && writer.statusCode != 201 {
				logger.Info("ResponseLog",
					zap.String("path", c.Request.RequestURI),
					zap.String("method", c.Request.Method),
					zap.Int("status", writer.statusCode),
					zap.String("user", fmt.Sprintf("%v/%v", "-", "-")),
					zap.String("requestBody", requestBody),
					zap.String("response", writer.body.String()),
				)
			}
		} else {
			// Log the response details using Zap
			if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" && writer.statusCode != 200 && writer.statusCode != 201 {
				logger.Info("ResponseLog",
					zap.String("path", c.Request.RequestURI),
					zap.String("method", c.Request.Method),
					zap.Int("status", writer.statusCode),
					zap.String("user", fmt.Sprintf("%v/%v/%v", userData.Email, userData.Username, userData.Role)),
					zap.String("requestBody", requestBody),
					zap.String("response", writer.body.String()),
				)
			}
		}
	}
}
