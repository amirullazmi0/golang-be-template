package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	SMTP     SMTPConfig
	Logger   LoggerConfig
}

type AppConfig struct {
	Name  string
	Env   string
	Port  string
	Debug bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

type JWTConfig struct {
	Secret      string
	ExpiredHour int
}

type CORSConfig struct {
	AllowedOrigins []string
}

type SMTPConfig struct {
	Email    string
	Password string
	Host     string
	Port     int
	FromName string
	FromEmail string
}

type LoggerConfig struct {
	LogToFile   bool
	LogFilePath string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Name:  viper.GetString("APP_NAME"),
			Env:   viper.GetString("APP_ENV"),
			Port:  viper.GetString("APP_PORT"),
			Debug: viper.GetBool("APP_DEBUG"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
			TimeZone: viper.GetString("DB_TIMEZONE"),
		},
		JWT: JWTConfig{
			Secret:      viper.GetString("JWT_SECRET"),
			ExpiredHour: viper.GetInt("JWT_EXPIRED_HOUR"),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
		},
		SMTP: SMTPConfig{
			Email:     viper.GetString("SMTP_EMAIL"),
			Password:  viper.GetString("SMTP_PASSWORD"),
			Host:      viper.GetString("SMTP_HOST"),
			Port:      viper.GetInt("SMTP_PORT"),
			FromName:  viper.GetString("SMTP_FROM_NAME"),
			FromEmail: viper.GetString("SMTP_FROM_EMAIL"),
		},
		Logger: LoggerConfig{
			LogToFile:   viper.GetBool("LOG_TO_FILE"),
			LogFilePath: viper.GetString("LOG_FILE_PATH"),
			MaxSize:     viper.GetInt("LOG_MAX_SIZE"),
			MaxBackups:  viper.GetInt("LOG_MAX_BACKUPS"),
			MaxAge:      viper.GetInt("LOG_MAX_AGE"),
			Compress:    viper.GetBool("LOG_COMPRESS"),
		},
	}

	return config, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode, c.TimeZone,
	)
}
