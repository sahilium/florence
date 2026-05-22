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
	Rooms            map[string]string

	TursoDatabaseURL string
	TursoAuthToken   string

	HNMinScore int
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
		Rooms: map[string]string{
			"hackernews": os.Getenv("ROOM_HN"),
			"sms":        os.Getenv("ROOM_SMS"),
			"lobsters":   os.Getenv("ROOM_LOBSTERS"),
			"github":     os.Getenv("ROOM_GITHUB"),
		},
		TursoDatabaseURL: getEnv("TURSO_DATABASE_URL", ""),
		TursoAuthToken:   getEnv("TURSO_AUTH_TOKEN", ""),
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
