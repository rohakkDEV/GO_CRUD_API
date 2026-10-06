package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	MongoDB    string
	ServerPort string
}

func Load() (Config, error) {
	// godotenv.Load reads .env and sets the values into the process env;
	// os.Getenv then reads them back.
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load env: %w", err)
	}

	mongoURI, err := extractEnv("MONGO_URI")
	if err != nil {
		return Config{}, err
	}

	mongoDB, err := extractEnv("MONGO_DB_NAME")
	if err != nil {
		return Config{}, err
	}

	port, err := extractEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	return Config{
		MongoURI:   mongoURI,
		MongoDB:    mongoDB,
		ServerPort: port,
	}, nil
}

func extractEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing required env var: %s", key)
	}
	return val, nil
}