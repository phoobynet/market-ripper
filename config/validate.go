package config

import (
	"errors"
	"log"
	"os"
)

// ValidateEnv validates the environment variables to ensure APCA_API_KEY_ID and APCA_API_SECRET_KEY are set
func ValidateEnv() {
	if key := os.Getenv("APCA_API_KEY_ID"); key == "" {
		log.Fatal(errors.New("APCA_API_KEY_ID is not set"))

	}

	if secret := os.Getenv("APCA_API_SECRET_KEY"); secret == "" {
		log.Fatal(errors.New("APCA_API_SECRET_KEY is not set"))
	}
}
