package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"

	kerr "github.com/krok-o/krok/errors"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/models"
)

const (
	defaultTokenExpiry        = time.Minute * 15
	defaultRefreshTokenExpiry = (time.Hour * 24) * 7
)

// TokenIssuerConfig contains the config for the TokenIssuer.
type TokenIssuerConfig struct {
	GlobalTokenKey string
}

// TokenIssuerDependencies contains the TokenIssuer dependencies.
type TokenIssuerDependencies struct {
	UserStore providers.UserStorer
	Clock     providers.Clock
}

// TokenIssuer represents the auth JWT token issuer.
type TokenIssuer struct {
	TokenIssuerConfig
	TokenIssuerDependencies
}

// NewTokenIssuer creates a new TokenIssuer.
func NewTokenIssuer(cfg TokenIssuerConfig, deps TokenIssuerDependencies) *TokenIssuer {
	return &TokenIssuer{TokenIssuerConfig: cfg, TokenIssuerDependencies: deps}
}

// Create creates a JWT access_token and refresh_token with the given user details.
// It will attempt to get or create the user in the database.
func (ti *TokenIssuer) Create(ctx context.Context, ud *models.UserAuthDetails) (*oauth2.Token, error) {
	user, err := ti.getOrCreateUser(ctx, ud)
	if err != nil {
		return nil, err
	}

	return ti.createToken(strconv.Itoa(user.ID))
}

func (ti *TokenIssuer) createToken(userID string) (*oauth2.Token, error) {
	now := ti.Clock.Now()

	// Create the new access token
	newAccessClaims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: now.Add(defaultTokenExpiry).Unix(),
		IssuedAt:  now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims).SignedString([]byte(ti.GlobalTokenKey))
	if err != nil {
		return nil, err
	}

	// Create the new refresh token
	newRefreshClaims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: now.Add(defaultRefreshTokenExpiry).Unix(),
		IssuedAt:  now.Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshClaims).SignedString([]byte(ti.GlobalTokenKey))
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		Expiry:       time.Unix(newAccessClaims.ExpiresAt, 0),
		RefreshToken: refreshToken,
	}, nil
}

func (ti *TokenIssuer) getOrCreateUser(ctx context.Context, ud *models.UserAuthDetails) (*models.User, error) {
	user, err := ti.UserStore.GetByEmail(ctx, ud.Email)
	if err != nil {
		var qe *kerr.QueryError
		if errors.As(err, &qe) && errors.Is(qe.Err, kerr.ErrNotFound) {
			// Not in the database, create them.
			dname := fmt.Sprintf("%s %s", ud.FirstName, ud.LastName)
			user, err = ti.UserStore.Create(ctx, &models.User{Email: ud.Email, DisplayName: dname})
			if err != nil {
				return nil, fmt.Errorf("create user: %w", err)
			}
			return user, nil
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

// Refresh refreshes the users JWT tokens.
func (ti *TokenIssuer) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	var refreshClaims jwt.StandardClaims
	// Parse & verify the refreshToken. Returns an error if the token has expired.
	if _, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ti.GlobalTokenKey), nil
	}); err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(refreshClaims.Subject)
	if err != nil {
		return nil, err
	}

	// Get the user. Allows us to check this person is still valid.
	user, err := ti.UserStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	newToken, err := ti.createToken(strconv.Itoa(user.ID))
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
