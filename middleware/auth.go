package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/DevilPepper/go-libs/auth"
	"github.com/DevilPepper/go-libs/auth/authiface"
)

type contextKey struct{}

var oidcClaimsKey = contextKey{}

func JWTAuthMiddleware() Middleware {
	verifier := auth.GetJWTVerifier()
	if verifier == nil {
		return NoOp()
	}
	return JWTAuthMiddlewareWithVerifier(verifier)
}

func JWTAuthMiddlewareWithVerifier(verifier authiface.JWTVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if verifier == nil {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := verifier.Verify(r.Context(), tokenString)
			if err != nil || token == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var claims auth.OIDCClaims
			if err := token.Claims(&claims); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), oidcClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetJWTClaims(ctx context.Context) *auth.OIDCClaims {
	claims, ok := ctx.Value(oidcClaimsKey).(auth.OIDCClaims)
	if !ok {
		return nil
	}
	return &claims
}
