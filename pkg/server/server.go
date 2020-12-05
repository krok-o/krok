package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"

	"github.com/krok-o/krok/pkg/krok"
)

const (
	api = "/rest/api/1"
)

// Config is the configuration of the bot's main cycle.
type Config struct {
	Port           string
	Hostname       string
	ServerKeyPath  string
	ServerCrtPath  string
	AutoTLS        bool
	CacheDir       string
	GlobalTokenKey string
}

// KrokServer is a server.
type KrokServer struct {
	Config
	Dependencies
}

// Dependencies defines needed dependencies for this bot.
type Dependencies struct {
	Logger zerolog.Logger
	Krok   krok.Handler
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
	// Setup Global Token Key
	if s.Config.GlobalTokenKey == "" {
		s.Logger.Info().Msg("Please set a global secret key... Randomly generating one for now...")
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		state := base64.StdEncoding.EncodeToString(b)
		s.Config.GlobalTokenKey = state
	}
	s.Dependencies.Logger.Info().Msg("Start listening...")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// The route will include a UUID with which we can use to identify this hook.
	e.POST("/hook/:id", s.Dependencies.Krok.HandleHooks(ctx))

	// Admin related actions
	//auth := e.Group(api+"/krok", middleware.JWT([]byte(s.Config.GlobalTokenKey)))
	//auth.GET("/hooks")

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
	// Start regular server
	return e.Start(hostPort)
}
