package cmd

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/krok-o/krok/pkg/krok"
	"github.com/krok-o/krok/pkg/krok/providers"
	"github.com/krok-o/krok/pkg/krok/providers/auth"
	"github.com/krok-o/krok/pkg/krok/providers/cache"
	"github.com/krok-o/krok/pkg/krok/providers/filevault"
	"github.com/krok-o/krok/pkg/krok/providers/handlers"
	"github.com/krok-o/krok/pkg/krok/providers/service"
	"github.com/krok-o/krok/pkg/krok/providers/vault"

	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
	"github.com/krok-o/krok/pkg/krok/providers/mailgun"
	"github.com/krok-o/krok/pkg/krok/providers/plugins"
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
		// 	OAuth
		oauthCfg auth.OAuthConfig
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
	flag.StringVar(&krokArgs.server.Port, "port", "9998", "--port 443")
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

	// Auth config
	flag.StringVar(&krokArgs.oauthCfg.GoogleClientID, "auth-google-client-id", "", "--auth-google-client-id my-client-id}")
	flag.StringVar(&krokArgs.oauthCfg.GoogleClientSecret, "auth-google-client-secret", "", "--auth-google-client-secret my-client-secret}")
	flag.StringVar(&krokArgs.oauthCfg.SessionSecret, "auth-session-secret", "1234", "--auth-session-secret my-secret}")
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
	a, _ := auth.NewKrokAuth(auth.AuthConfig{}, auth.AuthDependencies{
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

	userCache := cache.NewUserCache()
	go func() {
		interval := 1 * time.Minute
		for {
			userCache.ClearTTL()

			// nolint:gosimple
			select {
			case <-time.After(interval):
				log.Debug().Msg("Running user cache cleanup...")
			}
		}
	}()

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
		UserCache:  userCache,
	}
	tp, err := handlers.NewTokenProvider(handlers.Config{
		Hostname:       krokArgs.server.Hostname,
		GlobalTokenKey: krokArgs.server.GlobalTokenKey,
	}, handlerDeps)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create token handler.")
	}

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

	// ************************
	// Set up the server
	// ************************

	uuidGenerator := providers.NewUUIDGenerator()
	clock := providers.NewClock()

	oauthProvider := auth.NewOAuthProvider(krokArgs.oauthCfg, auth.OAuthProviderDependencies{
		Store:     userStore,
		UUID:      uuidGenerator,
		UserCache: userCache,
	})

	sv := server.NewKrokServer(krokArgs.server, server.Dependencies{
		Logger:         log,
		Krok:           krokHandler,
		CommandHandler: commandHandler,

		TokenProvider: tp,
		// RepositoryService
		RepositoryService: service.NewRepositoryService(service.RepositoryServiceConfig{
			Hostname: krokArgs.server.Hostname,
		}, service.RepositoryServiceDependencies{
			Logger: log,
			Storer: repoStore,
		}),
		// UserApiKeyService
		UserApiKeyService: service.NewUserAPIKeyService(service.UserAPIKeyServiceDependencies{
			Logger:        log,
			Storer:        apiKeyStore,
			Clock:         clock,
			Authenticator: authMatcher,
			UUID:          uuidGenerator,
		}),
		// 	AuthService
		AuthService: service.NewAuthService(oauthProvider),
	})

	// Run service & server
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return sv.Run(ctx)
	})

	g.Go(func() error {
		return sv.RunGRPC(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Err(err).Msg("Failed to run")
	}
}

// Execute runs the main krok command.
func Execute() error {
	return krokCmd.Execute()
}
