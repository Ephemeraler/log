package thlog

import (
	"fmt"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LOG_ENV_DEVELOPMENT = 1
	LOG_ENV_PRODUCTION  = 2
)

type Logger struct {
	environment int                // the environment for logger using.
	lumberjack  *lumberjack.Logger // for production to build logger.
	zap         *zap.Logger        // core logger.
}

// Close calls the underlying Core's Sync method in zap,
// flushing any buffered log entries. Applications should
// take care to call Sync before exiting. And in production,
// it will also close the lumberjack logger.
func (l *Logger) Close() {
	if l.environment == LOG_ENV_PRODUCTION {
		l.lumberjack.Close()
	}
	l.zap.Sync()
}

func (l *Logger) Record() *zap.Logger {
	return l.zap
}

func NewDevelopment() (*Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &Logger{
		environment: LOG_ENV_DEVELOPMENT,
		zap:         logger,
	}, nil
}

// NewProduction creates a new logger for production environment.
// The filename is the file to write logs to. The maxsize is the
// maximum size in megabytes of the log file before it gets
// rotated. The maxbackups is the maximum number of old log files
// to retain. The maxage is the maximum number of days to retain
// old log files based on the timestamp encoded in their filename.
func NewProduction(filename string, maxsize, maxbackups, maxage int) (*Logger, error) {
	// check the parameters.
	if filename == "" {
		return nil, fmt.Errorf("filename is empty")
	}

	if maxsize <= 0 {
		return nil, fmt.Errorf("maxsize is invalid")
	}

	if maxbackups <= 0 {
		return nil, fmt.Errorf("maxbackups is invalid")
	}

	if maxage <= 0 {
		return nil, fmt.Errorf("maxage is invalid")
	}

	// create the lumberjack logger.
	lumberjacklogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxsize, // megabytes
		MaxBackups: maxbackups,
		MaxAge:     maxage, //days
		Compress:   true,   // disabled by default
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339TimeEncoder
	fileEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(lumberjacklogger),
		zapcore.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller())

	return &Logger{
		environment: LOG_ENV_PRODUCTION,
		lumberjack:  lumberjacklogger,
		zap:         logger,
	}, nil
}
