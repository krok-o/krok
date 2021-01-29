package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/cache"
	"github.com/krok-o/krok/pkg/models"
)

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GlobalTokenKey     string
}

type OAuthProviderDependencies struct {
	Store     providers.UserStorer
	UUID      providers.UUIDGenerator
	UserCache *cache.UserCache
	Issuer    providers.TokenIssuer
}

// OAuthProvider is the OAuth provider.
type OAuthProvider struct {
	OAuthConfig
	OAuthProviderDependencies
	oauthCfg *oauth2.Config
}

func NewOAuthProvider(cfg OAuthConfig, deps OAuthProviderDependencies) *OAuthProvider {
	return &OAuthProvider{
		OAuthConfig:               cfg,
		OAuthProviderDependencies: deps,

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
func (op *OAuthProvider) GetAuthCodeURL(state string) string {
	return op.oauthCfg.AuthCodeURL(state, []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}...)
}

// Exchange exchanges the code returned from the OAuth2 authentication URL for a valid token.
// We attempt to get an internal user and create it if it doesn't exist.
func (op *OAuthProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := op.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	gu, err := op.getGoogleUser(token.AccessToken)
	if err != nil {
		return nil, err
	}

	var userID int
	if _, exists := op.UserCache.Has(gu.Email); !exists {
		user, err := op.Store.GetByEmail(ctx, gu.Email)
		if err != nil {
			var qe *kerr.QueryError
			if errors.As(err, &qe) && errors.Is(qe.Err, kerr.ErrNotFound) {
				// dname := fmt.Sprintf("%s %s", gu.FirstName, gu.LastName)
				user, err = op.Store.Create(ctx, &models.User{Email: gu.Email, DisplayName: "dname"})
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		userID = user.ID
	}
	op.UserCache.Add(gu.Email, userID)

	return op.Issuer.Create(providers.TokenProfile{UserID: strconv.Itoa(userID)})
}

type stateClaims struct {
	jwt.StandardClaims
	RedirectURL string `json:"redirect_url"`
}

// GenerateState generates the state nonce JWT with expiry.
func (op *OAuthProvider) GenerateState(redirectURL string) (string, error) {
	uuid, err := op.UUID.Generate()
	if err != nil {
		return "", err
	}

	claims := stateClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   uuid,
			ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		RedirectURL: redirectURL,
	}
	state, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(op.GlobalTokenKey))
	if err != nil {
		return "", err
	}

	return state, nil
}

// VerifyState verifies the state nonce JWT.
func (op *OAuthProvider) VerifyState(rawToken string) (string, error) {
	var claims stateClaims
	_, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(op.GlobalTokenKey), nil
	})
	if err != nil {
		return "", err
	}

	if err := claims.Valid(); err != nil {
		return "", err
	}

	return claims.RedirectURL, nil
}

type googleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

func (op *OAuthProvider) getGoogleUser(accessToken string) (*models.User, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("get user info: %s", err.Error())
	}
	defer response.Body.Close()

	var user *googleUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &models.User{
		Email: user.Email,
	}, nil
}
