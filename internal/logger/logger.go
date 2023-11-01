package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// Create a log file and configure log rotation
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Fatal("Failed to open log file:", err)
	}

	// Create a logrus hook for file logging
	logrusHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  logFile,
		logrus.ErrorLevel: logFile,
	}, &logrus.JSONFormatter{})

	// Add the file hook to the logger
	Logger.AddHook(logrusHook)

	// Set log level
	Logger.SetLevel(logrus.InfoLevel)

	// Configure console (standard output) logging
	Logger.SetOutput(os.Stdout)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Log the request and response
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		Logger.Printf("[%s] %s %s %v\n", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, latency)
	}
}
