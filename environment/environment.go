package environment

import (
	"os"
)

func GetBaseUrl() string {
	return os.Getenv("BASE_URL")
}
