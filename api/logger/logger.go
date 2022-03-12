// Package logger contains the logging functionality
// of the pmd-dx-api, onsisting of logger initialization,
// custom types and logging methods.
package logger

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var accessLogger *log.Logger
var errorLogger *log.Logger

var accessLogFile *lumberjack.Logger
var errorLogFile *lumberjack.Logger

// InitLogger opens all necessary log files and creates the log.Logger used by this package.
func InitLogger() error {
	// Get log path from environment
	logPath, ok := os.LookupEnv("LOG_PATH")
	if !ok {
		logPath = "logs"
	}
	// Use lumberjack instead of default logger for automated log rotation
	// Open the log files - error handling and flags are handled by the lumberjack package
	accessLogFile := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, "access.log"),
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
	}
	errorLogFile := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, "error.log"),
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
	}
	// Create the loggers
	accessLogger = log.New(accessLogFile, "", 0)
	errorLogger = log.New(errorLogFile, "", log.Ldate|log.Ltime)
	return nil
}

// CloseLogger closes the log files used by this package.
func CloseLogger() error {
	if accessLogFile == nil || errorLogFile == nil {
		return errors.New("no logging files to close")
	}
	err := accessLogFile.Close()
	if err != nil {
		return err
	}
	err = errorLogFile.Close()
	if err != nil {
		return err
	}
	return nil
}

// LogResponseRecorder is a custom http.ResponseWriter recording status and body size
// of a HTTP response for logging purposes.
type LogResponseRecorder struct {
	http.ResponseWriter
	Status int
	Size   int
}

// WriteHeader - implementation of http.ResponseWriter interface storing the status code.
func (l *LogResponseRecorder) WriteHeader(status int) {
	l.Status = status
	l.ResponseWriter.WriteHeader(status)
}

// Write - implementation of http.ResponseWriter interface storing the body size.
func (l *LogResponseRecorder) Write(b []byte) (int, error) {
	l.Size = len(b)
	return l.ResponseWriter.Write(b)
}

// LogRequest logs a HTTP request and the data of the ResponseRecorder to the accessLogger.
func LogRequest(request *http.Request, response LogResponseRecorder) error {
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

// CallerInformation represents information returned by
// runtime.Caller() and is used to pass caller information
// to the logger.
type CallerInformation struct {
	Pc   uintptr
	File string
	Line int
}

// String transforms the CallerInformation into a string representation.
func (c *CallerInformation) String() (string, error) {
	fileRegex, err := regexp.Compile(`\w*.go$`)
	if err != nil {
		return "", err
	}
	file := fileRegex.FindString(c.File)
	functionRegex, err := regexp.Compile(`\w*$`)
	if err != nil {
		return "", err
	}
	function := functionRegex.FindString(runtime.FuncForPC(c.Pc).Name())
	return fmt.Sprintf("%s(%s:%v)", function, file, c.Line), nil
}

// LogError logs an error to the errorLogger.
func LogError(err error, caller CallerInformation) error {
	if errorLogger == nil {
		return errors.New("error logger not initialized")
	}
	// Log the caller information and the error
	callerString, stringErr := caller.String()
	if stringErr != nil {
		return stringErr
	}
	errorLogger.Println(fmt.Sprintf("%s - %v", callerString, err.Error()))
	return nil
}
