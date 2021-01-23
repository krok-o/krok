package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/krok-o/krok/pkg/models"
	authv1 "github.com/krok-o/krok/proto/auth/v1"
)

// AuthServiceConfig represents the AuthService config.
type AuthServiceConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
}

// AuthService is the Authentication service.
type AuthService struct {
	googleOAuthConfig *oauth2.Config

	authv1.UnimplementedAuthServiceServer
}

// NewAuthService creates a AuthService.
func NewAuthService(cfg AuthServiceConfig) *AuthService {
	return &AuthService{
		googleOAuthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  "http://localhost:8081/auth.v1.AuthService/Callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// Login handles OAuth2 logins.
func (s *AuthService) Login(ctx context.Context, request *authv1.LoginRequest) (*empty.Empty, error) {
	if request.Provider == authv1.LoginProvider_GOOGLE {
		// TODO: Generate a random string in a cookie for state.
		url := s.googleOAuthConfig.AuthCodeURL("state")

		header := metadata.Pairs("Location", url)
		if err := grpc.SendHeader(ctx, header); err != nil {
			return nil, err
		}

		return &empty.Empty{}, nil
	}

	return nil, status.Error(codes.InvalidArgument, "invalid provider: "+request.Provider.String())
}

// Callback is the OAuth0 callback endpoint.
func (s *AuthService) Callback(ctx context.Context, request *authv1.CallbackRequest) (*empty.Empty, error) {
	// TODO: Design for more providers than Google.
	gu, err := s.getGoogleUser(ctx, request.Code)
	if err != nil {
		return nil, err
	}

	// Check if user/email exists.
	// Create if not exists.

	// TODO: DisplayName, FirstName, LastName.
	user := models.User{
		Email: gu.Email,
	}
	fmt.Println(user.Email)

	// Redirect with token

	return &empty.Empty{}, nil
}

type googleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

func (s *AuthService) getGoogleUser(ctx context.Context, code string) (*googleUser, error) {
	token, err := s.googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
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
