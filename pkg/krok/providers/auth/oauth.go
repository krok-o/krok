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
	"github.com/krok-o/krok/pkg/models"
)

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	SessionSecret      string
}

type OAuthProviderDependencies struct {
	Store providers.UserStorer
	UUID  providers.UUIDGenerator
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
		oauthCfg: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  "http://localhost:8081/auth.v1.AuthService/Callback?provider=GOOGLE",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (op *OAuthProvider) GetAuthCodeURL(state string) string {
	return op.oauthCfg.AuthCodeURL(state, []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}...)
}

func (op *OAuthProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := op.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	gu, err := op.getGoogleUser(token.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := op.Store.GetByEmail(ctx, gu.Email)
	if err != nil {
		var qe *kerr.QueryError
		if errors.As(err, &qe) && errors.Is(qe.Err, kerr.ErrNotFound) {
			dname := fmt.Sprintf("%s %s", gu.FirstName, gu.LastName)
			user, err = op.Store.Create(ctx, &models.User{Email: gu.Email, DisplayName: dname})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	claims := jwt.StandardClaims{
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(op.SessionSecret))
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Unix(claims.ExpiresAt, 0),
		TokenType:   "Bearer",
	}, nil
}

// GenerateState generates the state nonce JWT.
func (op *OAuthProvider) GenerateState() (string, error) {
	uuid, err := op.UUID.Generate()
	if err != nil {
		return "", err
	}

	claims := jwt.StandardClaims{
		Subject: uuid,
	}
	state, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(op.SessionSecret))
	if err != nil {
		return "", err
	}

	return state, nil
}

// VerifyState verifies the state nonce JWT.
func (op *OAuthProvider) VerifyState(rawToken string) error {
	var claims jwt.StandardClaims
	_, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(op.SessionSecret), nil
	})
	if err != nil {
		return err
	}

	if err := claims.Valid(); err != nil {
		return err
	}

	return nil
}

type googleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

func (op *OAuthProvider) getGoogleUser(accessToken string) (*googleUser, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("get user info: %s", err.Error())
	}
	defer response.Body.Close()

	var user *googleUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}
