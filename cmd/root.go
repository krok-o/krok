package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/krok-o/krok/pkg/krok"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/auth"
	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/github"
	"github.com/krok-o/krok/pkg/krok/providers/handlers"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/mailgun"
	"github.com/krok-o/krok/pkg/krok/providers/plugins"
	"github.com/krok-o/krok/pkg/krok/providers/vault"
	"github.com/krok-o/krok/pkg/models"
	"github.com/krok-o/krok/pkg/server"
)

var (
	krokCmd = &cobra.Command{
		Use:   "krok",
		Short: "Krok server",
		Run:   runKrokCmd,
	}
	krokArgs struct {
		devMode     bool
		debug       bool
		server      server.Config
		environment environment.Config
		store       livestore.Config
		plugins     plugins.Config
		email       mailgun.Config
		fileVault   filevault.Config
	}
)

func init() {
	flag := krokCmd.Flags()
	// Server Configs
	flag.BoolVar(&krokArgs.server.AutoTLS, "auto-tls", false, "--auto-tls")
	flag.BoolVar(&krokArgs.debug, "debug", false, "--debug")
	flag.StringVar(&krokArgs.server.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&krokArgs.server.ServerKeyPath, "server-key-path", "", "--server-key-path /home/user/.server/server.key")
	flag.StringVar(&krokArgs.server.ServerCrtPath, "server-crt-path", "", "--server-crt-path /home/user/.server/server.crt")
	flag.StringVar(&krokArgs.server.Proto, "proto", "http", "--proto http")
	flag.StringVar(&krokArgs.server.Hostname, "hostname", "localhost:9998", "--hostname localhost:9998")
	flag.StringVar(&krokArgs.server.GlobalTokenKey, "token", "", "--token <somerandomdata>")

	// Store config
	flag.StringVar(&krokArgs.store.Database, "krok-db-dbname", "krok", "--krok-db-dbname krok")
	flag.StringVar(&krokArgs.store.Username, "krok-db-username", "krok", "--krok-db-username krok")
	flag.StringVar(&krokArgs.store.Password, "krok-db-password", "password123", "--krok-db-password password123")
	flag.StringVar(&krokArgs.store.Hostname, "krok-db-hostname", "localhost:5432", "--krok-db-hostname localhost:5432")

	// Email
	flag.StringVar(&krokArgs.email.Domain, "email-domain", "", "--email-domain krok.com")
	flag.StringVar(&krokArgs.email.APIKey, "email-apikey", "", "--email-apikey ********")

	// Plugins
	flag.StringVar(&krokArgs.plugins.Location, "krok-plugin-location", "/tmp/krok/plugins", "--krok-plugin-location /tmp/krok/plugins")

	// VaultStorer config
	flag.StringVar(&krokArgs.fileVault.Location, "krok-file-vault-location", "/tmp/krok/vault", "--krok-file-vault-location /tmp/krok/vault")

	// OAuth
	flag.StringVar(&krokArgs.server.GoogleClientID, "google-client-id", "", "--google-client-id my-client-id}")
	flag.StringVar(&krokArgs.server.GoogleClientSecret, "google-client-secret", "", "--google-client-secret my-client-secret}")
}

// runKrokCmd builds up all the components and starts the krok server.
func runKrokCmd(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	out := zerolog.ConsoleWriter{
		Out: os.Stderr,
	}
	log := zerolog.New(out).With().
		Timestamp().
		Logger()

	if krokArgs.server.GoogleClientID == "" {
		log.Fatal().Msg("must provide --google-client-id flag")
	}
	if krokArgs.server.GoogleClientSecret == "" {
		log.Fatal().Msg("must provide --google-client-secret flag")
	}
	krokArgs.server.Addr = fmt.Sprintf("%s://%s", krokArgs.server.Proto, krokArgs.server.Hostname)

	// Setup Global Token Key
	if krokArgs.server.GlobalTokenKey == "" {
		log.Info().Msg("Please set a global secret key... Randomly generating one for now...")
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal().Msg("failed to generate global token key")
		}
		state := base64.StdEncoding.EncodeToString(b)
		krokArgs.server.GlobalTokenKey = state
	}

	// ************************
	// Set up db connection, vault and auth handlers.
	// ************************

	converter := environment.NewDockerConverter(environment.Config{}, environment.Dependencies{
		Logger: log,
	})
	connector := livestore.NewDatabaseConnector(krokArgs.store, livestore.Dependencies{
		Logger:    log,
		Converter: converter,
	})
	deps := livestore.Dependencies{
		Logger:    log,
		Converter: converter,
	}
	fv, _ := filevault.NewFileStorer(krokArgs.fileVault, filevault.Dependencies{
		Logger: log,
	})
	v, _ := vault.NewKrokVault(vault.Config{}, vault.Dependencies{
		Logger: log,
		Storer: fv,
	})
	a, _ := auth.NewRepositoryAuth(auth.RepositoryAuthConfig{}, auth.RepositoryAuthDependencies{
		Logger: log,
		Vault:  v,
	})

	// ************************
	// Set up stores
	// ************************

	repoStore := livestore.NewRepositoryStore(livestore.RepositoryDependencies{
		Dependencies: deps,
		Connector:    connector,
		Vault:        v,
		Auth:         a,
	})

	commandStore, err := livestore.NewCommandStore(livestore.CommandDependencies{
		Dependencies: deps,
		Connector:    connector,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to set up command store.")
	}

	apiKeyStore := livestore.NewAPIKeysStore(livestore.APIKeysDependencies{
		Dependencies: deps,
		Connector:    connector,
	})

	userStore := livestore.NewUserStore(livestore.UserDependencies{
		Dependencies: deps,
		Connector:    connector,
		APIKeys:      apiKeyStore,
	})

	// ************************
	// Set up platforms
	// ************************

	platformTokenProvider := auth.NewPlatformTokenProvider(auth.TokenProviderConfig{}, auth.TokenProviderDependencies{
		Logger: log,
		Vault:  v,
	})

	githubProvider := github.NewGithubPlatformProvider(github.Config{
		Hostname: krokArgs.server.Hostname,
	}, github.Dependencies{
		Logger:                log,
		AuthProvider:          a,
		PlatformTokenProvider: platformTokenProvider,
	})

	// ************************
	// Set up handlers
	// ************************
	authMatcher, err := auth.NewApiKeysProvider(auth.ApiKeysConfig{}, auth.ApiKeysDependencies{
		Logger:       log,
		ApiKeysStore: apiKeyStore,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new Api keys provider.")
	}
	handlerDeps := handlers.Dependencies{
		Logger:     log,
		UserStore:  userStore,
		ApiKeyAuth: authMatcher,
	}
	tp, err := handlers.NewTokenProvider(handlers.Config{
		Hostname:           krokArgs.server.Hostname,
		GlobalTokenKey:     krokArgs.server.GlobalTokenKey,
		GoogleClientID:     krokArgs.server.GoogleClientID,
		GoogleClientSecret: krokArgs.server.GoogleClientSecret,
	}, handlerDeps)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create token handler.")
	}

	platformProviders := make(map[int]providers.Platform)
	platformProviders[models.GITHUB] = githubProvider
	repoHandler, _ := handlers.NewRepositoryHandler(handlers.Config{
		Hostname:       krokArgs.server.Hostname,
		GlobalTokenKey: krokArgs.server.GlobalTokenKey,
	}, handlers.RepoHandlerDependencies{
		RepositoryStorer:  repoStore,
		TokenProvider:     tp,
		Logger:            log,
		PlatformProviders: platformProviders,
	})

	apiKeysHandler, _ := handlers.NewApiKeysHandler(handlers.Config{
		Hostname:       krokArgs.server.Hostname,
		GlobalTokenKey: krokArgs.server.GlobalTokenKey,
	}, handlers.ApiKeysHandlerDependencies{
		APIKeysStore:  apiKeyStore,
		TokenProvider: tp,
		Dependencies:  handlerDeps,
	})

	commandHandler, _ := handlers.NewCommandsHandler(handlers.Config{
		Hostname:       krokArgs.server.Hostname,
		GlobalTokenKey: krokArgs.server.GlobalTokenKey,
	}, handlers.CommandsHandlerDependencies{
		CommandStorer: commandStore,
		TokenProvider: tp,
		Logger:        log,
	})

	krokHandler := krok.NewHookHandler(krok.Config{}, krok.Dependencies{
		Logger: log,
	})

	tokenIssuer := auth.NewTokenIssuer(auth.TokenIssuerConfig{
		GlobalTokenKey: krokArgs.server.GlobalTokenKey,
	}, auth.TokenIssuerDependencies{
		UserStore: userStore,
		Clock:     providers.NewClock(),
	})

	uuidGenerator := providers.NewUUIDGenerator()
	oauthProvider := auth.NewOAuthAuthenticator(auth.OAuthAuthenticatorConfig{
		BaseURL:            krokArgs.server.Addr,
		GlobalTokenKey:     krokArgs.server.GlobalTokenKey,
		GoogleClientID:     krokArgs.server.GoogleClientID,
		GoogleClientSecret: krokArgs.server.GoogleClientSecret,
	}, auth.OAuthAuthenticatorDependencies{
		UUID:   uuidGenerator,
		Issuer: tokenIssuer,
		Clock:  providers.NewClock(),
	})

	authHandler := handlers.NewUserAuthHandler(handlers.UserAuthHandlerDeps{
		OAuthProvider: oauthProvider,
		TokenIssuer:   tokenIssuer,
		Logger:        log,
	})

	// ************************
	// Set up the server
	// ************************

	sv := server.NewKrokServer(krokArgs.server, server.Dependencies{
		Logger:            log,
		Krok:              krokHandler,
		CommandHandler:    commandHandler,
		RepositoryHandler: repoHandler,
		ApiKeyHandler:     apiKeysHandler,
		AuthHandler:       authHandler,
	})

	// Run service & server
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return sv.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Err(err).Msg("Failed to run")
	}
}

// Execute runs the main krok command.
func Execute() error {
	return krokCmd.Execute()
}
