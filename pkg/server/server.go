package server

import (
	"context"
	"errors"
	"net/http"

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
	CommandRunHandler      providers.CommandRunHandler
	RepositoryHandler      providers.RepositoryHandler
	APIKeyHandler          providers.APIKeysHandler
	AuthHandler            providers.AuthHandler
	TokenHandler           providers.TokenHandler
	VCSTokenHandler        providers.VCSTokenHandler
	SupportedPlatformList  providers.SupportedPlatformListHandler
	EventsHandler          providers.EventHandler
	VaultHandler           providers.VaultHandler
	UserHandler            providers.UserHandler
	ReadyHandler           providers.ReadyHandler
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

	// Healthz
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "alive")
	})
	e.GET("/readyz", s.ReadyHandler.Ready())

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
	auth.PUT("/command", s.Dependencies.CommandHandler.Upload())
	auth.POST("/command", s.Dependencies.CommandHandler.Create())
	auth.GET("/command/:id", s.Dependencies.CommandHandler.Get())
	auth.DELETE("/command/:id", s.Dependencies.CommandHandler.Delete())
	auth.POST("/commands", s.Dependencies.CommandHandler.List())
	auth.POST("/command/update", s.Dependencies.CommandHandler.Update())
	auth.POST("/command/add-command-rel-for-repository/:cmdid/:repoid", s.Dependencies.CommandHandler.AddCommandRelForRepository())
	auth.POST("/command/remove-command-rel-for-repository/:cmdid/:repoid", s.Dependencies.CommandHandler.RemoveCommandRelForRepository())
	auth.POST("/command/add-command-rel-for-platform/:cmdid/:pid", s.Dependencies.CommandHandler.AddCommandRelForPlatform())
	auth.POST("/command/remove-command-rel-for-platform/:cmdid/:pid", s.Dependencies.CommandHandler.RemoveCommandRelForPlatform())

	// command settings
	auth.GET("/command/settings/:id", s.Dependencies.CommandSettingsHandler.Get())
	auth.DELETE("/command/settings/:id", s.Dependencies.CommandSettingsHandler.Delete())
	auth.POST("/command/:id/settings", s.Dependencies.CommandSettingsHandler.List())
	auth.POST("/command/settings/update", s.Dependencies.CommandSettingsHandler.Update())
	auth.POST("/command/setting", s.Dependencies.CommandSettingsHandler.Create())

	// command runs
	auth.GET("/command/run/:id", s.Dependencies.CommandRunHandler.GetCommandRun())

	// api keys related actions
	auth.POST("/user/apikey/generate/:name", s.Dependencies.APIKeyHandler.Create())
	auth.DELETE("/user/apikey/delete/:keyid", s.Dependencies.APIKeyHandler.Delete())
	auth.GET("/user/apikeys", s.Dependencies.APIKeyHandler.List())
	auth.GET("/user/apikey/:keyid", s.Dependencies.APIKeyHandler.Get())

	// vcs token handler
	auth.POST("/vcs-token", s.Dependencies.VCSTokenHandler.Create())

	// events
	auth.POST("/events/:repoid", s.Dependencies.EventsHandler.List())
	auth.GET("/event/:id", s.Dependencies.EventsHandler.Get())

	// vault settings
	auth.POST("/vault/secret", s.Dependencies.VaultHandler.CreateSecret())
	auth.POST("/vault/secrets", s.Dependencies.VaultHandler.ListSecrets())
	auth.GET("/vault/secret/:name", s.Dependencies.VaultHandler.GetSecret())
	auth.POST("/vault/secret/update", s.Dependencies.VaultHandler.UpdateSecret())
	auth.DELETE("/vault/secret/:name", s.Dependencies.VaultHandler.DeleteSecret())

	// users
	auth.POST("/user", s.Dependencies.UserHandler.CreateUser())
	auth.POST("/users", s.Dependencies.UserHandler.ListUsers())
	auth.GET("/user/:id", s.Dependencies.UserHandler.GetUser())
	auth.POST("/user/update", s.Dependencies.UserHandler.UpdateUser())
	auth.DELETE("/user/:id", s.Dependencies.UserHandler.DeleteUser())

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
