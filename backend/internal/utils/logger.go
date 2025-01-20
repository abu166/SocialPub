package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// CustomFormatter extends logrus.JSONFormatter
type CustomFormatter struct {
	logrus.JSONFormatter
}

// Global logger instance
var Log *logrus.Logger

// LogConfig holds configuration for the logger
type LogConfig struct {
	LogLevel        string
	LogPath         string
	EnableConsole   bool
	EnableJSON      bool
	EnableRotation  bool
	MaxSize         int  // megabytes
	MaxBackups      int  // number of backups
	MaxAge          int  // days
	EnableCallerLog bool // Include caller information
}

// InitLogger initializes the logger with the given configuration
func InitLogger(config LogConfig) error {
	Log = logrus.New()

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(config.LogPath, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Set log level
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	Log.SetLevel(level)

	// Configure formatter
	if config.EnableJSON {
		Log.SetFormatter(&CustomFormatter{
			JSONFormatter: logrus.JSONFormatter{
				TimestampFormat:   time.RFC3339,
				DisableTimestamp:  false,
				DisableHTMLEscape: true,
				PrettyPrint:       false,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					filename := filepath.Base(f.File)
					return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
				},
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "timestamp",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
					logrus.FieldKeyFunc:  "caller",
				},
			},
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	}

	// Enable caller logging if configured
	if config.EnableCallerLog {
		Log.SetReportCaller(true)
	}

	// Setup log file
	currentTime := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("app-%s.log", currentTime)
	logFilePath := filepath.Join(config.LogPath, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Configure output writers
	var writers []io.Writer
	writers = append(writers, logFile)

	if config.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	// Set multi-writer if needed
	if len(writers) > 1 {
		Log.SetOutput(io.MultiWriter(writers...))
	} else {
		Log.SetOutput(writers[0])
	}

	return nil
}

// LoggerMiddleware creates a middleware that logs HTTP requests
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(startTime)

		// Log the request details
		Log.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     rw.statusCode,
			"duration":   duration.String(),
			"ip":         r.RemoteAddr,
			"user_agent": r.UserAgent(),
		}).Info("HTTP Request")
	})
}

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Helper logging functions
func LogError(err error, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["error"] = err
	Log.WithFields(fields).Error(err.Error())
}

func LogInfo(message string, fields logrus.Fields) {
	Log.WithFields(fields).Info(message)
}

func LogWarn(message string, fields logrus.Fields) {
	Log.WithFields(fields).Warn(message)
}

func LogDebug(message string, fields logrus.Fields) {
	Log.WithFields(fields).Debug(message)
}

// LogRequest logs the details of an HTTP request
func LogRequest(r *http.Request, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["method"] = r.Method
	fields["path"] = r.URL.Path
	fields["remote_addr"] = r.RemoteAddr
	fields["user_agent"] = r.UserAgent()

	Log.WithFields(fields).Info("Request received")
}

// LogResponse logs the details of an HTTP response
func LogResponse(statusCode int, duration time.Duration, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["status_code"] = statusCode
	fields["duration"] = duration.String()

	Log.WithFields(fields).Info("Response sent")
}

// RotateLogs checks if log rotation is needed and performs rotation if necessary
func RotateLogs(config LogConfig) error {
	if !config.EnableRotation {
		return nil
	}

	logDir := config.LogPath
	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %v", err)
	}

	var totalSize int64
	var oldestFile os.FileInfo
	oldestTime := time.Now()

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".log" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		totalSize += info.Size()

		if info.ModTime().Before(oldestTime) {
			oldestTime = info.ModTime()
			oldestFile = info
		}
	}

	// Check if total size exceeds max size
	if totalSize > int64(config.MaxSize)*1024*1024 && oldestFile != nil {
		oldPath := filepath.Join(logDir, oldestFile.Name())
		backupPath := filepath.Join(logDir, fmt.Sprintf("%s.%s.backup",
			oldestFile.Name(), time.Now().Format("20060102150405")))

		if err := os.Rename(oldPath, backupPath); err != nil {
			return fmt.Errorf("failed to rotate log file: %v", err)
		}

		// Remove old backups if exceeding MaxBackups
		if config.MaxBackups > 0 {
			removeOldBackups(logDir, config.MaxBackups)
		}
	}

	return nil
}

// removeOldBackups removes old backup files exceeding the maximum count
func removeOldBackups(logDir string, maxBackups int) error {
	files, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	var backups []os.FileInfo
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".backup" {
			info, err := file.Info()
			if err != nil {
				continue
			}
			backups = append(backups, info)
		}
	}

	// Sort backups by modification time (oldest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime().Before(backups[j].ModTime())
	})

	// Remove excess backups
	for i := 0; i < len(backups)-maxBackups; i++ {
		path := filepath.Join(logDir, backups[i].Name())
		if err := os.Remove(path); err != nil {
			Log.WithError(err).Warn("Failed to remove old backup")
		}
	}

	return nil
}
