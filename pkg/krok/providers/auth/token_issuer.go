package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const (
	defaultTokenExpiry        = time.Minute * 15
	defaultRefreshTokenExpiry = (time.Hour * 24) * 7
)

type TokenIssuerConfig struct {
	GlobalTokenKey string
}

type TokenIssuer struct {
	cfg TokenIssuerConfig
}

func NewTokenIssuer(cfg TokenIssuerConfig) *TokenIssuer {
	return &TokenIssuer{cfg: cfg}
}

func (ti *TokenIssuer) Create(token providers.TokenProfile) (*oauth2.Token, error) {
	newAccessClaims := jwt.StandardClaims{
		Subject:   token.UserID,
		ExpiresAt: time.Now().Add(defaultTokenExpiry).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims).SignedString([]byte(ti.cfg.GlobalTokenKey))
	if err != nil {
		return nil, err
	}

	newRefreshClaims := jwt.StandardClaims{
		Subject:   token.UserID,
		ExpiresAt: time.Now().Add(defaultRefreshTokenExpiry).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshClaims).SignedString([]byte(ti.cfg.GlobalTokenKey))
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

func (ti *TokenIssuer) Refresh(refreshToken string) (*oauth2.Token, error) {
	var refreshClaims jwt.StandardClaims
	if _, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ti.cfg.GlobalTokenKey), nil
	}); err != nil {
		return nil, err
	}

	newToken, err := ti.Create(providers.TokenProfile{UserID: refreshClaims.Subject})
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
