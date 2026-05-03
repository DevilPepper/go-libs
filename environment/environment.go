package environment

import (
	"os"
)

func GetBaseUrl() string {
	return os.Getenv("BASE_URL")
}

func GetDomain() string {
	if IS_DEV {
		domain := os.Getenv("DOMAIN")
		if domain == "" {
			domain = "localhost"
		}
		return domain
	}
	return ""
}
