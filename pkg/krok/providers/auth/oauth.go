package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

// OAuthAuthenticatorConfig contains the config for the OAuthAuthenticator.
type OAuthAuthenticatorConfig struct {
	BaseURL            string
	GoogleClientID     string
	GoogleClientSecret string
	GlobalTokenKey     string
}

// OAuthAuthenticatorDependencies contains the dependencies for the OAuthAuthenticator.
type OAuthAuthenticatorDependencies struct {
	UUID      providers.UUIDGenerator
	Clock     providers.Clock
	Issuer    providers.TokenIssuer
	UserStore providers.UserStorer
}

// OAuthAuthenticator is the OAuthAuthenticator that uses OAuth2 to authenticate the user.
type OAuthAuthenticator struct {
	OAuthAuthenticatorConfig
	OAuthAuthenticatorDependencies

	oauthCfg *oauth2.Config
}

// NewOAuthAuthenticator creates a new OAuthAuthenticator.
func NewOAuthAuthenticator(cfg OAuthAuthenticatorConfig, deps OAuthAuthenticatorDependencies) *OAuthAuthenticator {
	return &OAuthAuthenticator{
		OAuthAuthenticatorConfig:       cfg,
		OAuthAuthenticatorDependencies: deps,

		// TODO: Support multiple providers.
		oauthCfg: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  fmt.Sprintf("%s/auth/callback", cfg.BaseURL),
			Scopes:       scopes,
			Endpoint:     google.Endpoint,
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

	ud, err := op.getGoogleUser(token.AccessToken)
	if err != nil {
		return nil, err
	}

	// TODO: In future we may want a constraint where users must be created first, rather than automatically.
	user, err := op.getOrCreateUser(ctx, ud)
	if err != nil {
		return nil, err
	}

	return op.Issuer.Create(user)
}

func (op *OAuthAuthenticator) getOrCreateUser(ctx context.Context, ud *models.UserAuthDetails) (*models.User, error) {
	user, err := op.UserStore.GetByEmail(ctx, ud.Email)
	if err != nil {
		var qe *kerr.QueryError
		if errors.As(err, &qe) && errors.Is(qe.Err, kerr.ErrNotFound) {
			// Not in the database, create them.
			dname := fmt.Sprintf("%s %s", ud.FirstName, ud.LastName)
			user, err = op.UserStore.Create(ctx, &models.User{Email: ud.Email, DisplayName: dname})
			if err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}
			return user, nil
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
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
	if _, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(op.GlobalTokenKey), nil
	}); err != nil {
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
