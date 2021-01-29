package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"

	"github.com/krok-o/krok/pkg/krok"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/handlers"
	krokmiddleware "github.com/krok-o/krok/pkg/server/middleware"
)

const (
	api = "/rest/api/1"
)

// Config is the configuration of the server
type Config struct {
	Port               string
	Hostname           string
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
	Logger            zerolog.Logger
	Krok              krok.Handler
	CommandHandler    providers.CommandHandler
	RepositoryHandler providers.RepositoryHandler
	ApiKeyHandler     providers.ApiKeysHandler
	AuthHandler       providers.AuthHandler

	// TokenProvider providers.TokenProvider
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
	e.GET("/auth/login", s.AuthHandler.Login())
	e.GET("/auth/callback", s.AuthHandler.Callback())

	// Routes
	// This is the general format of a hook callback url for a repository.
	// @rid repository id
	// @vid vcs id
	e.POST("/hooks/:rid/:vid/callback", s.Dependencies.Krok.HandleHooks(ctx))

	userAuthMiddleware := krokmiddleware.JWTAuthentication(&krokmiddleware.JWTAuthConfig{
		CookieName:     handlers.AccessTokenCookie,
		GlobalTokenKey: s.GlobalTokenKey,
	})
	auth := e.Group(api+"/krok", userAuthMiddleware)

	auth.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	// Repository related actions.
	auth.POST("/repository", s.Dependencies.RepositoryHandler.CreateRepository())
	auth.GET("/repository/:id", s.Dependencies.RepositoryHandler.GetRepository())
	auth.DELETE("/repository/:id", s.Dependencies.RepositoryHandler.DeleteRepository())
	auth.POST("/repositories", s.Dependencies.RepositoryHandler.ListRepositories())
	auth.POST("/repository/update", s.Dependencies.RepositoryHandler.UpdateRepository())

	// command related actions.
	auth.GET("/command/:id", s.Dependencies.CommandHandler.GetCommand())
	auth.DELETE("/command/:id", s.Dependencies.CommandHandler.DeleteCommand())
	auth.POST("/commands", s.Dependencies.CommandHandler.ListCommands())
	auth.POST("/command/update", s.Dependencies.CommandHandler.UpdateCommand())
	auth.POST("/command/add-command-rel-for-repository/:cmdid/:repoid", s.Dependencies.CommandHandler.AddCommandRelForRepository())
	auth.POST("/command/remove-command-rel-for-repository/:cmdid/:repoid", s.Dependencies.CommandHandler.RemoveCommandRelForRepository())

	// api keys related actions
	auth.POST("/user/:uid/apikey/generate/:name", s.Dependencies.ApiKeyHandler.CreateApiKeyPair())
	auth.DELETE("/user/:uid/apikey/delete/:keyid", s.Dependencies.ApiKeyHandler.DeleteApiKeyPair())
	auth.POST("/user/:uid/apikeys", s.Dependencies.ApiKeyHandler.ListApiKeyPairs())
	auth.GET("/user/:uid/apikey/:keyid", s.Dependencies.ApiKeyHandler.GetApiKeyPair())

	hostPort := fmt.Sprintf("%s:%s", s.Config.Hostname, s.Config.Port)

	// Start TLS with certificate paths
	if len(s.Config.ServerKeyPath) > 0 && len(s.Config.ServerCrtPath) > 0 {
		e.Pre(middleware.HTTPSRedirect())
		return e.StartTLS(hostPort, s.Config.ServerCrtPath, s.Config.ServerKeyPath)
	}

	// Start Auto TLS server
	if s.Config.AutoTLS {
		if len(s.Config.CacheDir) < 1 {
			return errors.New("cache dir must be provided if autoTLS is enabled")
		}
		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.Cache = autocert.DirCache(s.Config.CacheDir)
		return e.StartAutoTLS(hostPort)
	}

	go func() {
		<-ctx.Done()
		if err := e.Shutdown(ctx); err != nil {
			s.Logger.Debug().Err(err).Msg("Failed to shutdown the server nicely.")
		}
	}()

	// Start regular server
	return e.Start(hostPort)
}
