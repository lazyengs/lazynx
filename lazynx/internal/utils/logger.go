package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupFileLogger creates a zap logger that writes to a file
func SetupFileLogger(logFile string, verbose bool) (*zap.SugaredLogger, error) {
	// Ensure log directory exists
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// Create or append to log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Configure log level
	level := zapcore.InfoLevel
	if verbose {
		level = zapcore.DebugLevel
	}

	// Create file writer
	fileWriter := zapcore.AddSync(file)

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create core with file output
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		level,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar(), nil
}

// GetDefaultLogFile returns the default log file path for lazynx
func GetDefaultLogFile() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".lazynx", "logs", "lazynx.log")
}
