package authiface

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCProvider interface {
	VerifierContext(ctx context.Context, config *oidc.Config) *oidc.IDTokenVerifier
}

type JWTVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

type OIDC interface {
	NewProvider(issuerUrl string) (OIDCProvider, error)
	NewVerifier(oidcProvider OIDCProvider, clientId string) JWTVerifier
}
