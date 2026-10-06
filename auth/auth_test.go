//go:build !prod && !qa

package auth

import (
	"testing"

	"github.com/DevilPepper/go-libs/generated/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetJWTVerifierWithFactory_NoIssuer(t *testing.T) {
	t.Setenv("ISSUER_URL", "")
	t.Setenv("CLIENT_ID", "test-client")

	mockFactory := mocks.NewMockOIDC(t)
	verifier, err := GetJWTVerifierWithFactory(mockFactory)
	assert.Nil(t, verifier)
	assert.NoError(t, err)
	mockFactory.AssertNotCalled(t, "NewProvider")
}

func TestGetJWTVerifierWithFactory_NoClientID(t *testing.T) {
	t.Setenv("ISSUER_URL", "https://example.com")
	t.Setenv("CLIENT_ID", "")

	mockFactory := mocks.NewMockOIDC(t)
	mockProvider := mocks.NewMockOIDCProvider(t)
	mockFactory.On("NewProvider", "https://example.com").Return(mockProvider, nil)

	verifier, err := GetJWTVerifierWithFactory(mockFactory)
	assert.Nil(t, verifier)
	assert.Error(t, err)
	mockFactory.AssertExpectations(t)
}

func TestGetJWTVerifierWithFactory_ProviderError(t *testing.T) {
	t.Setenv("ISSUER_URL", "https://example.com")
	t.Setenv("CLIENT_ID", "test-client")

	mockFactory := mocks.NewMockOIDC(t)
	mockFactory.On("NewProvider", "https://example.com").Return(nil, assert.AnError)

	verifier, err := GetJWTVerifierWithFactory(mockFactory)
	assert.Nil(t, verifier)
	assert.ErrorIs(t, err, assert.AnError)
	mockFactory.AssertExpectations(t)
}

func TestGetJWTVerifierWithFactory_Success(t *testing.T) {
	t.Setenv("ISSUER_URL", "https://example.com")
	t.Setenv("CLIENT_ID", "test-client")

	mockFactory := mocks.NewMockOIDC(t)
	mockVerifier := mocks.NewMockJWTVerifier(t)
	mockFactory.On("NewProvider", "https://example.com").Return(nil, nil)
	mockFactory.On("NewVerifier", mock.Anything, "test-client").Return(mockVerifier)

	verifier, err := GetJWTVerifierWithFactory(mockFactory)
	assert.NoError(t, err)
	assert.Equal(t, mockVerifier, verifier)
	mockFactory.AssertExpectations(t)
}
