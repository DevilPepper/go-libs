package auth

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/DevilPepper/go-libs/auth/authiface"
	"github.com/DevilPepper/go-libs/environment"
	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	jwtVerifier authiface.JWTVerifier
	once        sync.Once
)

type OIDC struct {
	ctx context.Context
}

type OIDCClaims struct {
	Email         string   `json:"email"`
	Name          string   `json:"name"`
	Groups        []string `json:"groups"`
	EmailVerified bool     `json:"email_verified"`
	Subject       string   `json:"sub"`
}

func (o *OIDC) NewProvider(issuerUrl string) (authiface.OIDCProvider, error) {
	return oidc.NewProvider(o.ctx, issuerUrl)
}
func (o *OIDC) NewVerifier(oidcProvider authiface.OIDCProvider, clientId string) authiface.JWTVerifier {
	return oidcProvider.VerifierContext(o.ctx, &oidc.Config{ClientID: clientId})
}

func GetJWTVerifier() authiface.JWTVerifier {
	if jwtVerifier == nil {
		once.Do(func() {
			verifier, err := GetJWTVerifierWithFactory(&OIDC{ctx: context.Background()})
			if err != nil {
				panic(err)
			}
			jwtVerifier = verifier
		})
	}
	return jwtVerifier
}

func GetJWTVerifierWithFactory(oidcFactory authiface.OIDC) (authiface.JWTVerifier, error) {
	issuerURL := os.Getenv("ISSUER_URL")
	if issuerURL == "" {
		if !environment.IS_DEV {
			return nil, errors.New("ISSUER_URL environment variable is required")
		} else {
			return nil, nil
		}
	}

	provider, err := oidcFactory.NewProvider(issuerURL)
	if err != nil {
		return nil, err
	}
	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("CLIENT_ID environment variable is required")
	}
	return oidcFactory.NewVerifier(provider, clientId), nil
}
