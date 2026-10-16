package auth

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/browser"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

//go:embed callback.html
var callbackHTML string

func Login(ctx context.Context, issuerURL string, clientID string, clientSecret string, app string, additionalScopes ...string) error {
	scopes := append(
		[]string{
			"openid",
			"profile",
			"email",
		},
		additionalScopes...,
	)

	token, err := loadToken(app, clientID)
	if err == nil && (token.RefreshToken != "" || token.Valid()) {
		hasRefresh := "❌"
		if token.RefreshToken != "" {
			hasRefresh = "✅"
		}
		log.Warn(
			"Already logged in",
			"refreshToken", hasRefresh,
			"accessTokenExpires", token.Expiry.In(time.Local).Format("03:04 PM"),
		)
		return nil
	}
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return err
	}

	verifier := oauth2.GenerateVerifier()
	oauthCfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:18000",
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	listener, err := net.Listen("tcp", "127.0.0.1:18000")
	if err != nil {
		return err
	}

	done := make(chan *oauth2.Token, 1)
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := oauthCfg.Exchange(ctx, r.URL.Query().Get("code"), oauth2.VerifierOption(verifier))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		done <- token
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(callbackHTML))
	})}
	go server.Serve(listener)

	browser.OpenURL(oauthCfg.AuthCodeURL("", oauth2.S256ChallengeOption(verifier), oauth2.AccessTypeOffline))
	token = <-done
	server.Shutdown(context.Background())

	if token.RefreshToken == "" {
		log.Infof("Logged in until %s", token.Expiry.In(time.Local).Format("03:04 PM"))
	}
	return saveToken(app, clientID, token)
}

type OIDCEndpoints struct {
	RevocationEndpoint string `json:"revocation_endpoint"`
}

func Logout(ctx context.Context, issuerURL string, clientID string, clientSecret string, app string) error {
	token, err := loadToken(app, clientID)
	if err != nil {
		return err
	}
	if token.RefreshToken != "" {
		log.Warn("No refresh token to revoke")
	} else {
		provider, err := oidc.NewProvider(ctx, issuerURL)
		if err != nil {
			return err
		}
		var endpoints OIDCEndpoints
		if err := provider.Claims(&endpoints); err != nil {
			return err
		}
		if endpoints.RevocationEndpoint == "" {
			log.Warn("IdP doesn't expose a revocation endpoint.")
		} else {
			form := url.Values{
				"token":           {token.RefreshToken},
				"token_type_hint": {"refresh_token"},
				"client_id":       {clientID},
			}
			if clientSecret != "" {
				form.Set("client_secret", clientSecret)
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoints.RevocationEndpoint, strings.NewReader(form.Encode()))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			err = doRequest(req)
			if err != nil {
				return fmt.Errorf(
					"Failed to revoke token: %w\n\n"+
						"Your session is still active. Try again later, or run:\n"+
						"  keyring del %s %s\n"+
						"to clear your local credentials anyway",
					err,
					app,
					clientID,
				)
			}
		}
	}
	return keyring.Delete(app, clientID)
}

func doRequest(req *http.Request) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := strings.TrimSpace(string(body))
		return fmt.Errorf("%s\n%s", resp.Status, bodyStr)
	}

	return nil
}

func AuthClient(ctx context.Context, issuerURL string, clientID string, clientSecret string, app string) *http.Client {
	provider, _ := oidc.NewProvider(ctx, issuerURL)
	oauthCfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
	}
	src := oauth2.ReuseTokenSource(nil, &refresher{config: oauthCfg, app: app})
	return oauth2.NewClient(ctx, src)
}

type refresher struct {
	config oauth2.Config
	app    string
}

func (r *refresher) Token() (*oauth2.Token, error) {
	token, err := loadToken(r.app, r.config.ClientID)
	if err != nil {
		return nil, err
	}
	newToken, err := r.config.TokenSource(context.Background(), token).Token()
	if err != nil {
		if strings.Contains(err.Error(), "token expired") && token.RefreshToken == "" {
			return nil, fmt.Errorf("Session expired. Please log in again")
		}
		return nil, err
	}
	if (newToken.RefreshToken != "" && newToken.RefreshToken != token.RefreshToken) || newToken.AccessToken != token.AccessToken {
		saveToken(r.app, r.config.ClientID, newToken)
	}
	return newToken, nil
}

func saveToken(app string, clientId string, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return keyring.Set(app, clientId, string(data))
}

func loadToken(app string, clientId string) (*oauth2.Token, error) {
	data, err := keyring.Get(app, clientId)
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, fmt.Errorf("Failed to retrieve auth token from keyring. Please logging in")
	}
	var token oauth2.Token
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
