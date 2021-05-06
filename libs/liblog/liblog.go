package liblog

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger обертка над логрусом.
type Logger struct {
	*logrus.Entry
}

// WithRequestFields sets logger fields with data from request.
func WithRequestFields(logger Logger, r *http.Request) Logger {
	return Logger{logger.WithFields(logrus.Fields{
		"client_ip":      r.Header.Get("X-Forwarded-For"),
		"user_agent":     r.Header.Get("User-Agent"),
		"session_id":     r.Header.Get("X-Authentication-Token"),
		"request_id":     r.Header.Get("X-Request-Id"),
		"request_path":   r.URL.Path,
		"request_query":  r.URL.RawQuery,
		"request":        r.URL.RequestURI(),
		"request_method": r.Method,
	})}
}

// LoggerConfig знает конфигурацию компонента логгирования.
type LoggerConfig struct {
	Level           string `env:"LOGGER_LEVEL"`
	Output          string `env:"LOGGER_OUTPUT"`
	Formatter       string `env:"LOGGER_FORMATTER"`
	HumanReadable   bool   `env:"LOGGER_HUMAN_READABLE"`
	TimeStampFormat string `env:"LOGGER_TIMESTAMP_FORMAT"`
	Caller          bool   `env:"LOGGER_CALLER"`
}

// NewLogger создает новый инстанс логгера из конфигурации
func NewLogger(config LoggerConfig) (Logger, error) {
	l := logrus.New()

	logger := logrus.NewEntry(l)
	if level, err := logrus.ParseLevel(config.Level); err == nil {
		logger.Logger.SetLevel(level)
	} else {
		logger.Fatal(err)
	}

	var formatter logrus.Formatter
	switch config.Formatter {
	case "json":
		formatter = &logrus.JSONFormatter{
			TimestampFormat: config.TimeStampFormat,
			PrettyPrint:     config.HumanReadable,
		}
	case "text":
		formatter = &logrus.TextFormatter{
			TimestampFormat: config.TimeStampFormat,
		}
	default:
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		}
	}
	logger.Logger.SetFormatter(formatter)

	var output io.Writer
	switch config.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		f, err := os.OpenFile(config.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return Logger{}, err
		}
		output = f
	}
	logger.Logger.SetOutput(output)

	logger.Logger.SetReportCaller(config.Caller)

	hostname, err := os.Hostname()
	if err != nil {
		logger.Warn(err)
	}

	logger = logger.WithFields(logrus.Fields{
		"pid":      os.Getpid(),
		"ppid":     os.Getppid(),
		"hostname": hostname,
	})
	return Logger{logger}, nil
}
