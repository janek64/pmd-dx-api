// Package logger contains the logging functionality of
// the pmd-dx-api, consisting of logger initialization,
// custom types and logging functions.
package logger

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var accessLogger *log.Logger

var accessLogFile *lumberjack.Logger

// InitLogger opens all necessary log files and creates the log.Logger used by this package.
func InitLogger() error {
	// Get log path from environment
	logPath, ok := os.LookupEnv("LOG_PATH")
	if !ok {
		logPath = "logs"
	}
	// Use lumberjack instead of default logger for automated log rotation
	// Open the log file - error handling and flags are handled by the lumberjack package
	accessLogFile := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, "access.log"),
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
	}
	// Create the logger
	accessLogger = log.New(accessLogFile, "", 0)
	return nil
}

// CloseLogger closes the log files used by this package.
func CloseLogger() error {
	if accessLogFile == nil {
		return errors.New("no logging files to close")
	}
	err := accessLogFile.Close()
	if err != nil {
		return err
	}
	return nil
}

// ResponseRecorder is a custom http.ResponseWriter recording status and body size
// of a HTTP response for logging purposes.
type ResponseRecorder struct {
	http.ResponseWriter
	Status int
	Size   int
}

// WriteHeader - implementation of http.ResponseWriter interface storing the status code.
func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write - implementation of http.ResponseWriter interface storing the body size.
func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Size = len(b)
	return r.ResponseWriter.Write(b)
}

// LogRequest logs a HTTP request and the data of the ResponseRecorder to the accessLogger.
func LogRequest(request *http.Request, response ResponseRecorder) error {
	if accessLogger == nil {
		return errors.New("access logger not initialized")
	}
	// Logging in "Combined Log Format" without referrer
	t := time.Now()
	accessLogger.Printf("%s - - [%s] \"%s %s %s\" %v %v \"%s\"\n",
		request.RemoteAddr,
		t.Format("02/Jan/2006:15:04:05 -0700"),
		request.Method,
		request.URL,
		request.Proto,
		response.Status,
		response.Size,
		request.UserAgent(),
	)
	return nil
}
