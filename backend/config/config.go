package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Environment     string
	MongoURI        string
	MongoDatabase   string
	JWTSecret       string
	JWTExpiration   string
	MaxUploadSize   int64
	UploadDir       string
	GridFSBucket    string
}

var AppConfig *Config

func LoadConfig() {
	// Загрузить .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	maxSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	AppConfig = &Config{
		Port:          getEnv("PORT", "8080"),
		Environment:   getEnv("ENV", "development"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "cinema_booking_db"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret-key"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
		MaxUploadSize: maxSize,
		UploadDir:     getEnv("UPLOAD_DIR", "../uploads"),
		GridFSBucket:  getEnv("GRIDFS_BUCKET", "cinema_files"),
	}

	log.Println("✅ Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}