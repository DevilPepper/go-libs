package environment

import (
	"os"
)

func IsDev() bool {
	environment := os.Getenv("ENVIRONMENT")
	return environment == "" || environment == "dev"
}

func GetBaseUrl() string {
	return os.Getenv("BASE_URL")
}
