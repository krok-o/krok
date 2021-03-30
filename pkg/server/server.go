package server

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"

	"github.com/krok-o/krok/pkg/krok/providers"
)

const (
	api = "/rest/api/1"
)

// Config is the configuration of the server
type Config struct {
	Proto              string
	Hostname           string
	HookBase           string
	Addr               string
	ServerKeyPath      string
	ServerCrtPath      string
	AutoTLS            bool
	CacheDir           string
	GlobalTokenKey     string
	GoogleClientID     string
	GoogleClientSecret string
}

// KrokServer is a server.
type KrokServer struct {
	Config
	Dependencies
}

// Dependencies defines needed dependencies for the krok server.
type Dependencies struct {
	Logger                 zerolog.Logger
	HookHandler            providers.HookHandler
	UserMiddleware         providers.UserMiddleware
	CommandHandler         providers.CommandHandler
	CommandSettingsHandler providers.CommandSettingsHandler
	RepositoryHandler      providers.RepositoryHandler
	APIKeyHandler          providers.APIKeysHandler
	AuthHandler            providers.AuthHandler
	TokenHandler           providers.TokenHandler
	VCSTokenHandler        providers.VCSTokenHandler
	UserTokenHandler       providers.UserTokenHandler
	SupportedPlatformList  providers.SupportedPlatformListHandler
}

// Server defines a server which runs and accepts requests.
type Server interface {
	Run(context.Context) error
}

// NewKrokServer creates a new krok server.
func NewKrokServer(cfg Config, deps Dependencies) *KrokServer {
	return &KrokServer{Config: cfg, Dependencies: deps}
}

// Run starts up listening.
func (s *KrokServer) Run(ctx context.Context) error {
	s.Dependencies.Logger.Info().Msg("Start listening...")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	// Public endpoints for authentication.
	e.POST("/auth/refresh", s.AuthHandler.Refresh())
	e.GET("/auth/login", s.AuthHandler.OAuthLogin())
	e.GET("/auth/callback", s.AuthHandler.OAuthCallback())
	e.GET("/supported-platforms", s.SupportedPlatformList.ListSupportedPlatforms())

	// Routes
	// This is the general format of a hook callback url for a repository.
	// @rid repository id
	// @vid vcs id
	e.POST(api+"/hooks/:rid/:vid/callback", s.Dependencies.HookHandler.HandleHooks())
	e.POST(api+"/get-token", s.Dependencies.TokenHandler.TokenHandler())
	// Admin related actions

	auth := e.Group(api+"/krok", s.Dependencies.UserMiddleware.JWT())

	// Repository related actions.
	auth.POST("/repository", s.Dependencies.RepositoryHandler.Create())
	auth.GET("/repository/:id", s.Dependencies.RepositoryHandler.Get())
	auth.DELETE("/repository/:id", s.Dependencies.RepositoryHandler.Delete())
	auth.POST("/repositories", s.Dependencies.RepositoryHandler.List())
	auth.POST("/repository/update", s.Dependencies.RepositoryHandler.Update())

	// command related actions.
	auth.GET("/command/:id", s.Dependencies.CommandHandler.Get())
	auth.DELETE("/command/:id", s.Dependencies.CommandHandler.Delete())
	auth.POST("/commands", s.Dependencies.CommandHandler.List())
	auth.POST("/command/update", s.Dependencies.CommandHandler.Update())
	auth.POST("/command/add-command-rel-for-repository/:cmdid/:repoid", s.Dependencies.CommandHandler.AddCommandRelForRepository())
	auth.POST("/command/remove-command-rel-for-repository/:cmdid/:repoid", s.Dependencies.CommandHandler.RemoveCommandRelForRepository())

	// command settings
	auth.GET("/command/settings/:id", s.Dependencies.CommandSettingsHandler.Get())
	auth.DELETE("/command/settings/:id", s.Dependencies.CommandSettingsHandler.Delete())
	auth.POST("/command/:id/settings", s.Dependencies.CommandSettingsHandler.List())
	auth.POST("/command/settings/update", s.Dependencies.CommandSettingsHandler.Update())
	auth.POST("/command/setting", s.Dependencies.CommandSettingsHandler.Create())

	// api keys related actions
	auth.POST("/user/apikey/generate/:name", s.Dependencies.APIKeyHandler.Create())
	auth.DELETE("/user/apikey/delete/:keyid", s.Dependencies.APIKeyHandler.Delete())
	auth.GET("/user/apikey", s.Dependencies.APIKeyHandler.List())
	auth.GET("/user/apikey/:keyid", s.Dependencies.APIKeyHandler.Get())

	// user personal token (api token)
	auth.POST("/user/token/generate", s.Dependencies.UserTokenHandler.Generate())

	// vcs token handler
	auth.POST("/vcs-token", s.Dependencies.VCSTokenHandler.Create())

	// Start TLS with certificate paths
	if len(s.Config.ServerKeyPath) > 0 && len(s.Config.ServerCrtPath) > 0 {
		e.Pre(middleware.HTTPSRedirect())
		return e.StartTLS(s.Config.Hostname, s.Config.ServerCrtPath, s.Config.ServerKeyPath)
	}

	// Start Auto TLS server
	if s.Config.AutoTLS {
		if len(s.Config.CacheDir) < 1 {
			return errors.New("cache dir must be provided if autoTLS is enabled")
		}
		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.Cache = autocert.DirCache(s.Config.CacheDir)
		return e.StartAutoTLS(s.Config.Hostname)
	}

	go func() {
		<-ctx.Done()
		if err := e.Shutdown(ctx); err != nil {
			s.Logger.Debug().Err(err).Msg("Failed to shutdown the server nicely.")
		}
	}()

	// Start regular server
	return e.Start(s.Config.Hostname)
}
