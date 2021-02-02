package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

// OAuthAuthenticatorConfig contains the config for the OAuthAuthenticator.
type OAuthAuthenticatorConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GlobalTokenKey     string
}

// OAuthAuthenticatorDependencies contains the dependencies for the OAuthAuthenticator.
type OAuthAuthenticatorDependencies struct {
	UUID   providers.UUIDGenerator
	Clock  providers.Clock
	Issuer providers.UserTokenIssuer
}

// OAuthAuthenticator is the OAuthAuthenticator that uses OAuth2 to authenticate the user.
type OAuthAuthenticator struct {
	OAuthAuthenticatorConfig
	OAuthAuthenticatorDependencies

	// TODO: Have a map of configs and support multiple providers.
	oauthCfg *oauth2.Config
}

// NewOAuthAuthenticator creates a new OAuthAuthenticator.
func NewOAuthAuthenticator(cfg OAuthAuthenticatorConfig, deps OAuthAuthenticatorDependencies) *OAuthAuthenticator {
	return &OAuthAuthenticator{
		OAuthAuthenticatorConfig:       cfg,
		OAuthAuthenticatorDependencies: deps,

		// For now, just support Google.
		oauthCfg: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  "http://localhost:9998/auth/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// GetAuthCodeURL gets the OAuth2 authentication URL.
func (op *OAuthAuthenticator) GetAuthCodeURL(state string) string {
	return op.oauthCfg.AuthCodeURL(state, []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}...)
}

// Exchange exchanges the code returned from the OAuth2 authentication URL for a valid token.
// We then call the TokenIssuer to get an internal JWT.
func (op *OAuthAuthenticator) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := op.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	userAuthDetails, err := op.getGoogleUser(token.AccessToken)
	if err != nil {
		return nil, err
	}

	return op.Issuer.Create(ctx, userAuthDetails)
}

// stateClaims are used when creating a temporary JWT state nonce that has an expiry.
// This is used for CSRF protection when logging in via an OAuth2 provider.
type stateClaims struct {
	jwt.StandardClaims
	RedirectURL string `json:"redirect_url"`
}

// GenerateState generates the state nonce JWT with expiry.
func (op *OAuthAuthenticator) GenerateState(redirectURL string) (string, error) {
	uuid, err := op.UUID.Generate()
	if err != nil {
		return "", fmt.Errorf("uuid generate: %w", err)
	}

	now := op.Clock.Now()
	claims := stateClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   uuid,
			ExpiresAt: now.Add(time.Minute * 2).Unix(),
			IssuedAt:  now.Unix(),
		},
		RedirectURL: redirectURL,
	}
	state, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(op.GlobalTokenKey))
	if err != nil {
		return "", fmt.Errorf("new signed token: %w", err)
	}

	return state, nil
}

// VerifyState verifies the state nonce JWT.
func (op *OAuthAuthenticator) VerifyState(rawToken string) (string, error) {
	var claims stateClaims
	_, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(op.GlobalTokenKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	return claims.RedirectURL, nil
}

type googleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

func (op *OAuthAuthenticator) getGoogleUser(accessToken string) (*models.UserAuthDetails, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("get user info: %s", err.Error())
	}
	defer response.Body.Close()

	var user *googleUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &models.UserAuthDetails{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
