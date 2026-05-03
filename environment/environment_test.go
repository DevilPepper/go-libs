package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBaseUrl(t *testing.T) {
	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv("BASE_URL", "https://example.com")
		assert.Equal(t, "https://example.com", GetBaseUrl())
	})

	t.Run("returns empty string when not set", func(t *testing.T) {
		os.Unsetenv("BASE_URL")
		assert.Empty(t, GetBaseUrl())
	})
}
