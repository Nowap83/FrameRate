package config

import (
	"fmt"
	"os"
)

func ValidateEnvironment() error {
	required := []string{
		"JWT_SECRET",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"TMDB_API_KEY",
		"TMDB_BASE_URL",
		"TMDB_IMAGE_BASE_URL",
		"RESEND_API_KEY",
		"FRONTEND_URL",
	}

	var missing []string
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}
	return nil
}
