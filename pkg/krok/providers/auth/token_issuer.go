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

type TokenIssuerConfig struct {
	GlobalTokenKey string
}

type TokenIssuerDependencies struct {
	UserCache providers.UserCache
	UserStore providers.UserStorer
}

type TokenIssuer struct {
	TokenIssuerConfig
	TokenIssuerDependencies
}

func NewTokenIssuer(cfg TokenIssuerConfig, deps TokenIssuerDependencies) *TokenIssuer {
	return &TokenIssuer{TokenIssuerConfig: cfg, TokenIssuerDependencies: deps}
}

func (ti *TokenIssuer) Create(ctx context.Context, ud *models.UserAuthDetails) (*oauth2.Token, error) {
	user, err := ti.getOrCreateUser(ctx, ud)
	if err != nil {
		return nil, err
	}

	return ti.createToken(strconv.Itoa(user.ID))
}

func (ti *TokenIssuer) createToken(userID string) (*oauth2.Token, error) {
	// Create the new access token
	newAccessClaims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(defaultTokenExpiry).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims).SignedString([]byte(ti.GlobalTokenKey))
	if err != nil {
		return nil, err
	}

	// Create the new refresh token
	newRefreshClaims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(defaultRefreshTokenExpiry).Unix(),
		IssuedAt:  time.Now().Unix(),
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

func (ti *TokenIssuer) getOrCreateUser(ctx context.Context, ud *models.UserAuthDetails) (user *models.User, err error) {
	defer func() {
		if user != nil {
			ti.UserCache.Add(ud.Email, user.ID)
		}
	}()

	// Check the cache for the user.
	if u, exists := ti.UserCache.Has(ud.Email); exists {
		return u.User, nil
	}

	// Not in the cache, check the database.
	u, err := ti.UserStore.GetByEmail(ctx, ud.Email)
	if err != nil {
		var qe *kerr.QueryError
		if errors.As(err, &qe) && errors.Is(qe.Err, kerr.ErrNotFound) {
			// Not in the database, createToken them.
			dname := fmt.Sprintf("%s %s", ud.FirstName, ud.LastName)
			user, err = ti.UserStore.Create(ctx, &models.User{Email: ud.Email, DisplayName: dname})
			if err != nil {
				return nil, err
			}
			return user, nil
		} else {
			return nil, err
		}
	}

	return u, nil
}

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
	ti.UserCache.Add(user.Email, user.ID)

	newToken, err := ti.createToken(strconv.Itoa(user.ID))
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
