package service

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/krok-o/krok/pkg/krok/providers"
	authv1 "github.com/krok-o/krok/proto/auth/v1"
)

// AuthService is the Authentication service.
type AuthService struct {
	oauthProvider providers.OAuthProvider

	authv1.UnimplementedAuthServiceServer
}

// NewAuthService creates a AuthService.
func NewAuthService(oauthProvider providers.OAuthProvider) *AuthService {
	return &AuthService{oauthProvider: oauthProvider}
}

// Login handles OAuth2 logins.
func (s *AuthService) Login(ctx context.Context, request *authv1.LoginRequest) (*empty.Empty, error) {
	state, err := s.oauthProvider.GenerateState(request.RedirectURL)
	if err != nil {
		return nil, err
	}

	url := s.oauthProvider.GetAuthCodeURL(state)

	// Redirect to the provided auth url.
	header := metadata.Pairs("Location", url)
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Callback is the OAuth0 callback endpoint.
func (s *AuthService) Callback(ctx context.Context, request *authv1.CallbackRequest) (*authv1.CallbackResponse, error) {
	redirectURL, err := s.oauthProvider.VerifyState(request.State)
	if err != nil {
		return nil, err
	}

	token, err := s.oauthProvider.Exchange(ctx, request.Code)
	if err != nil {
		return nil, err
	}

	header := metadata.Pairs(
		"Location", redirectURL,
		"Set-Cookie-Token", token.AccessToken,
	)
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, err
	}

	fmt.Println(token.AccessToken)

	return &authv1.CallbackResponse{
		AccessToken: token.AccessToken,
	}, nil
}
