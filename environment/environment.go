package environment

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
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

func AssertVariables(required ...string) {
	var missing []string
	for _, variable := range required {
		if os.Getenv(variable) == "" {
			missing = append(missing, variable)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing environment variable: \n - %s", strings.Join(missing, "\n - "))
	}
}
