//go:build prod

package auth

import (
	"testing"

	"github.com/DevilPepper/go-libs/generated/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetJWTVerifierWithFactory_NoIssuer_prod(t *testing.T) {
	t.Setenv("ISSUER_URL", "")
	t.Setenv("CLIENT_ID", "test-client")

	mockFactory := mocks.NewMockOIDC(t)
	verifier, err := GetJWTVerifierWithFactory(mockFactory)
	assert.Nil(t, verifier)
	assert.Error(t, err)
	mockFactory.AssertNotCalled(t, "NewProvider")
}
