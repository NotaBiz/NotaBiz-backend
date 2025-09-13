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

// responseWriter is a wrapper around gin.ResponseWriter that captures
// the response body and status code for logging purposes.
type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

// Write overrides the default gin.ResponseWriter Write method
// to also store the response body in memory for later logging.
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// WriteHeader overrides the default gin.ResponseWriter WriteHeader method
// to capture the HTTP status code.
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// SetupLogger configures and initializes a Zap logger with Lumberjack
// for log rotation. It uses JSON encoding in production and a
// development-friendly format in development environments.
//
// Logs are stored in files named with the current date under the configured path.
//
// Parameters:
//   - configData: application configuration containing logger settings.
//
// Returns:
//   - *zap.Logger: configured Zap logger instance
//   - error: error if logger setup fails
func SetupLogger(configData *model.ConfigData) (logging *zap.Logger, err error) {
	// Setup Lumberjack for log rotation
	logger := lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/log_%s.log", configData.LoggerConfig.Path, time.Now().Format("2006-01-02")),
		MaxSize:    configData.LoggerConfig.MaxSize,
		MaxBackups: configData.LoggerConfig.MaxBackups,
		MaxAge:     configData.LoggerConfig.MaxAge,
		Compress:   configData.LoggerConfig.Compress,
	}

	// Configure logging level and encoder
	var level zapcore.Level
	var encoderConfig zapcore.EncoderConfig
	if configData.AppConfig.Environment == "development" {
		level = zapcore.DebugLevel
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		level = zapcore.InfoLevel
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	// Create Zap logger
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(&logger), level)
	logging = zap.New(core, zap.AddCaller())

	return logging, nil
}

// ResponseLogger is a Gin middleware that logs request and response data
// for non-GET and non-OPTIONS requests where the response status code
// is not 200 or 201.
//
// Logged fields include:
//   - path: request URI
//   - method: HTTP method
//   - status: HTTP status code
//   - user: extracted from JWT claims (email/role) if available
//   - requestBody: body of the incoming request
//   - response: response body sent back to the client
//
// Parameters:
//   - logger: Zap logger instance for structured logging.
//
// Returns:
//   - gin.HandlerFunc: middleware function for Gin.
func ResponseLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Failed to read request body", zap.Error(err))
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Replace writer to capture response
		writer := &responseWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = writer

		// Process request
		c.Next()

		// Prepare log fields
		requestBody := string(bodyBytes)
		user, _ := c.Get("user")
		userData, ok := user.(*model.JwtClaims)

		// Log only if request is not GET/OPTIONS and response status is not 200/201
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" &&
			writer.statusCode != 200 && writer.statusCode != 201 {
			if ok {
				// Log with user information
				logger.Info("ResponseLog",
					zap.String("path", c.Request.RequestURI),
					zap.String("method", c.Request.Method),
					zap.Int("status", writer.statusCode),
					zap.String("user", fmt.Sprintf("%v/%v/%v", userData.Email, userData.Role)),
					zap.String("requestBody", requestBody),
					zap.String("response", writer.body.String()),
				)
			} else {
				// Log with anonymous user
				logger.Info("ResponseLog",
					zap.String("path", c.Request.RequestURI),
					zap.String("method", c.Request.Method),
					zap.Int("status", writer.statusCode),
					zap.String("user", "-/-"),
					zap.String("requestBody", requestBody),
					zap.String("response", writer.body.String()),
				)
			}
		}
	}
}
