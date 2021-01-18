package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/krok-o/krok/pkg/krok"
	"github.com/krok-o/krok/pkg/krok/providers"
	grpcmiddleware "github.com/krok-o/krok/pkg/server/middleware"
	repov1 "github.com/krok-o/krok/proto/repository/v1"
)

const (
	api = "/rest/api/1"
)

// Config is the configuration of the server
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

// Dependencies defines needed dependencies for the krok server.
type Dependencies struct {
	Logger         zerolog.Logger
	Krok           krok.Handler
	CommandHandler providers.CommandHandler
	ApiKeyHandler  providers.ApiKeysHandler

	TokenProvider     providers.TokenProvider
	RepositoryService repov1.RepositoryServiceServer
}

// Server defines a server which runs and accepts requests.
type Server interface {
	Run(context.Context) error
	RunGRPC(context.Context) error
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

	// Repository related actions.
	auth := e.Group(api+"/krok", middleware.JWT([]byte(s.Config.GlobalTokenKey)))

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
		e.Shutdown(ctx)
	}()

	// Start regular server
	return e.Start(hostPort)
}

// RunGRPC runs grpc and grpc-gateway.
func (s *KrokServer) RunGRPC(ctx context.Context) error {
	// TODO: Use SSL/TLS
	gs := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.JwtAuthInterceptor(s.TokenProvider)),
	)

	repov1.RegisterRepositoryServiceServer(gs, s.RepositoryService)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := repov1.RegisterRepositoryServiceHandlerFromEndpoint(ctx, mux, ":9090", opts); err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return gs.Serve(listener)
	})

	g.Go(func() error {
		return http.ListenAndServe(":8081", mux)
	})

	return g.Wait()
}
