package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MatrixHomeserver string
	MatrixToken      string
	MatrixRoomID     string
	HNMinScore       int
}

func Load() Config {
	_ = godotenv.Load()

	score, err := strconv.Atoi(getEnv("HN_MIN_SCORE", "100"))
	if err != nil {
		log.Fatal(err)
	}

	return Config{
		MatrixHomeserver: getEnv("MATRIX_HOMESERVER", ""),
		MatrixToken:      getEnv("MATRIX_ACCESS_TOKEN", ""),
		MatrixRoomID:     getEnv("MATRIX_ROOM_ID", ""),
		HNMinScore:       score,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
