package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port              string
	RedisAddr         string
	MongoURI          string
	RabbitMQURL       string
	FirebaseProjectID string
	DevMode           bool
	AdminUserIDs      map[string]bool
	DefaultShowtimeID string
}

func Load() *Config {
	loadDotEnv(".env")

	adminIDs := make(map[string]bool)
	for _, id := range strings.Split(os.Getenv("ADMIN_USER_IDS"), ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			adminIDs[id] = true
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	showtimeID := os.Getenv("DEFAULT_SHOWTIME_ID")
	if showtimeID == "" {
		showtimeID = "show_001"
	}

	return &Config{
		Port:              port,
		RedisAddr:         redisAddr,
		MongoURI:          mongoURI,
		RabbitMQURL:       rabbitURL,
		FirebaseProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
		DevMode:           os.Getenv("DEV_MODE") == "true",
		AdminUserIDs:      adminIDs,
		DefaultShowtimeID: showtimeID,
	}
}

func (c *Config) UserRole(userID string) string {
	if c.AdminUserIDs[userID] {
		return "ADMIN"
	}
	return "USER"
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
