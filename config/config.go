package config

import (
	"doc_bot/libs/liblog"
	"errors"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Database знает конфигурацию подключения к базе данных.
type Database struct {
	// SQL-Dialect
	Dialect         string        `env:"DB_DIALECT"`
	User            string        `env:"DB_USER"`
	Password        string        `env:"DB_PASSWORD"`
	Name            string        `env:"DB_NAME"`
	Host            string        `env:"DB_HOST"`
	Port            int           `env:"DB_PORT"`
	SSLMode         string        `env:"DB_SSL_MODE"`
	Schema          string        `env:"DB_SCHEMA"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `end:"DB_MAX_IDLE_CONNS"`
}

// Validate проверяет валидность конфигурации подключения к базе данных.
func (d Database) Validate() error {
	if d.Port < 0 {
		return errors.New("Invalid database port (DB_PORT)")
	}
	if d.User == "" {
		return errors.New("Invalid database user (DB_USER)")
	}
	if d.Password == "" {
		return errors.New("Invalid database password (DB_PASSWORD)")
	}
	if d.Name == "" {
		return errors.New("Invalid database name (DB_NAME)")
	}
	if d.Host == "" {
		return errors.New("Invalid database host (DB_HOST)")
	}
	if d.SSLMode == "" {
		return errors.New("Invalid database SSL mode (DB_SSL_MODE)")
	}
	if d.Schema == "" {
		return errors.New("Invalid database schema (DB_SCHEMA)")
	}

	return nil
}

// Transport знает конфигурацию транспортного уровня сервиса.
type Transport struct {
	Address string `env:"TRANSPORT_ADDRESS"`
	Port    int    `env:"TRANSPORT_PORT"`
}

// Validate проверяет валидность конфигурации транспортного уровня.
func (t Transport) Validate() error {
	if t.Port < 0 {
		return errors.New("Invalid transport port (TRANSPORT_PORT)")
	}
	return nil
}

// Telegram знает конфигурацию бота телеграм.
type Telegram struct {
	Token string `env:"TELEGRAM_TOKEN"`
}

// Validate проверяет валидность конфигурации бота телеграм.
func (t Telegram) Validate() error {
	if len(t.Token) == 0 {
		return errors.New("Invalid telegram token (TELEGRAM_TOKEN)")
	}
	return nil
}

// Config структура конфигурации сервиса.
type Config struct {
	Transport Transport
	Logger    liblog.LoggerConfig
	Telegram  Telegram
	Database  Database
}

// LoadConfig загружает конфигурацию из переменных среды ENV.
func LoadConfig() (*Config, error) {
	cfg := Default()
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Default возвращает конфигурацию по умолчанию
func Default() *Config {
	return &Config{
		Transport: Transport{
			Port: 8080,
		},
		Database: Database{
			Dialect:         "postgres",
			ConnMaxLifetime: 10 * time.Minute,
			MaxIdleConns:    25,
			MaxOpenConns:    25,
		},
		Logger: liblog.LoggerConfig{
			Level:     "warn",
			Output:    "stdout",
			Formatter: "json",
		},
	}
}

// Validate проверяет валидность конфигурации сервиса.
func (c Config) Validate() error {
	var err error

	err = c.Transport.Validate()
	if err != nil {
		return err
	}

	return nil
}
