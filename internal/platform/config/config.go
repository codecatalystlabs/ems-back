package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	Kafka KafkaConfig
	JWT   JWTConfig
	Log   LogConfig
}

type AppConfig struct {
	Name            string
	Env             string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Addr       string
	Password   string
	DB         int
	DefaultTTL time.Duration
}

type KafkaConfig struct {
	Brokers               []string
	ClientID              string
	GroupID               string
	TopicUserCreated      string
	TopicIncidentCreated  string
	TopicDispatchAssigned string
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type LogConfig struct {
	Level string
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		App: AppConfig{
			Name:            getenv("APP_NAME", "dispatch-backend"),
			Env:             getenv("APP_ENV", "development"),
			Port:            getenv("APP_PORT", "8080"),
			ReadTimeout:     mustDuration("APP_READ_TIMEOUT", "15s"),
			WriteTimeout:    mustDuration("APP_WRITE_TIMEOUT", "15s"),
			ShutdownTimeout: mustDuration("APP_SHUTDOWN_TIMEOUT", "10s"),
		},
		DB: DBConfig{
			Host:            getenv("DB_HOST", "localhost"),
			Port:            getenv("DB_PORT", "5432"),
			User:            getenv("DB_USER", "postgres"),
			Password:        getenv("DB_PASSWORD", "pwaiswa"),
			Name:            getenv("DB_NAME", "dispatch_db"),
			SSLMode:         getenv("DB_SSLMODE", "disable"),
			MaxOpenConns:    mustInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    mustInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: mustDuration("DB_CONN_MAX_LIFETIME", "30m"),
		},
		Redis: RedisConfig{
			Addr:       getenv("REDIS_ADDR", "localhost:6379"),
			Password:   getenv("REDIS_PASSWORD", ""),
			DB:         mustInt("REDIS_DB", 0),
			DefaultTTL: mustDuration("REDIS_DEFAULT_TTL", "5m"),
		},
		Kafka: KafkaConfig{
			Brokers:               strings.Split(getenv("KAFKA_BROKERS", "localhost:9092"), ","),
			ClientID:              getenv("KAFKA_CLIENT_ID", "dispatch-backend"),
			GroupID:               getenv("KAFKA_GROUP_ID", "dispatch-backend-group"),
			TopicUserCreated:      getenv("KAFKA_TOPIC_USER_CREATED", "user.created"),
			TopicIncidentCreated:  getenv("KAFKA_TOPIC_INCIDENT_CREATED", "incident.created"),
			TopicDispatchAssigned: getenv("KAFKA_TOPIC_DISPATCH_ASSIGNED", "dispatch.assigned"),
		},
		JWT: JWTConfig{
			Secret:     getenv("JWT_SECRET", "change-me"),
			Issuer:     getenv("JWT_ISSUER", "dispatch-backend"),
			AccessTTL:  mustDuration("JWT_ACCESS_TTL", "15m"),
			RefreshTTL: mustDuration("JWT_REFRESH_TTL", "168h"),
		},
		Log: LogConfig{Level: getenv("LOG_LEVEL", "debug")},
	}

	if cfg.JWT.Secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func mustDuration(key, fallback string) time.Duration {
	v := getenv(key, fallback)
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(err)
	}
	return d
}

func mustInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return n
}
