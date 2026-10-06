package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"unsafe"

	"github.com/DevilPepper/go-libs/auth"
	"github.com/DevilPepper/go-libs/generated/mocks"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJWTAuthMiddlewareWithVerifier_NoVerifier_PassesThrough(t *testing.T) {
	middleware := JWTAuthMiddlewareWithVerifier(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	rec := httptest.NewRecorder()

	middleware(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestJWTAuthMiddlewareWithVerifier_NoAuthHeader_Returns401(t *testing.T) {
	mockVerifier := mocks.NewMockJWTVerifier(t)
	middleware := JWTAuthMiddlewareWithVerifier(mockVerifier)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
	mockVerifier.AssertNotCalled(t, "Verify")
}

func TestJWTAuthMiddlewareWithVerifier_WrongAuthType_Returns401(t *testing.T) {
	mockVerifier := mocks.NewMockJWTVerifier(t)
	middleware := JWTAuthMiddlewareWithVerifier(mockVerifier)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	rec := httptest.NewRecorder()

	middleware(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
	mockVerifier.AssertNotCalled(t, "Verify")
}

func TestGetJWTClaims_FromContext(t *testing.T) {
	claims := auth.OIDCClaims{
		Email:         "test@example.com",
		Name:          "Test User",
		Groups:        []string{"admin", "users"},
		EmailVerified: true,
		Subject:       "user123",
	}

	ctx := context.WithValue(context.Background(), oidcClaimsKey, claims)
	retrieved := GetJWTClaims(ctx)

	if retrieved == nil {
		t.Fatal("expected claims, got nil")
	}
	if retrieved.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", retrieved.Email)
	}
	if retrieved.Subject != "user123" {
		t.Errorf("expected subject user123, got %s", retrieved.Subject)
	}
}

func TestJWTAuthMiddlewareWithVerifier_InvalidToken_Returns401(t *testing.T) {
	mockVerifier := mocks.NewMockJWTVerifier(t)
	mockVerifier.On("Verify", mock.Anything, mock.Anything).Return(nil, errors.New("invalid"))
	middleware := JWTAuthMiddlewareWithVerifier(mockVerifier)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	middleware(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestJWTAuthMiddlewareWithVerifier_HappyPath(t *testing.T) {
	jwtPayload := `{"sub":"user123","email":"test@example.com","name":"Test User","email_verified":true,"groups":["admin","users"],"iss":"https://example.com","aud":"test-client","exp":9999999999,"iat":1000000000}`

	mockVerifier := mocks.NewMockJWTVerifier(t)

	// Using reflection because `claim`, the thing we care about, is not exposed publicly
	idToken := &oidc.IDToken{}
	claimsField := reflect.ValueOf(idToken).Elem().FieldByName("claims")
	*(*[]byte)(unsafe.Pointer(claimsField.UnsafeAddr())) = []byte(jwtPayload)

	mockVerifier.On("Verify", mock.Anything, mock.Anything).Return(idToken, nil)
	middleware := JWTAuthMiddlewareWithVerifier(mockVerifier)

	var retrievedClaims *auth.OIDCClaims
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retrievedClaims = GetJWTClaims(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	middleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, retrievedClaims)
	assert.Equal(t, "user123", retrievedClaims.Subject)
	assert.Equal(t, "test@example.com", retrievedClaims.Email)
	assert.Equal(t, "Test User", retrievedClaims.Name)
	assert.True(t, retrievedClaims.EmailVerified)
	assert.Equal(t, []string{"admin", "users"}, retrievedClaims.Groups)

	mockVerifier.AssertExpectations(t)
}
