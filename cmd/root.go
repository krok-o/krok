package cmd

import (
	"context"
	"os"

	"github.com/krok-o/krok/pkg/krok"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

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

	// Store config
	flag.StringVar(&krokArgs.store.Database, "krok-db-dbname", "krok", "--krok-db-dbname krok")
	flag.StringVar(&krokArgs.store.Username, "krok-db-username", "krok", "--krok-db-username krok")
	flag.StringVar(&krokArgs.store.Password, "krok-db-password", "krok", "--krok-db-password password123")
	flag.StringVar(&krokArgs.store.Hostname, "krok-db-hostname", "", "--krok-db-hostname krok-db")

	// Email
	flag.StringVar(&krokArgs.email.Domain, "email-domain", "", "--email-domain krok.com")
	flag.StringVar(&krokArgs.email.APIKey, "email-apikey", "", "--email-apikey ********")

	// Plugins
	flag.StringVar(&krokArgs.plugins.Location, "krok-plugin-location", "/tmp/krok/plugins", "--krok-plugin-location /tmp/krok/plugins")
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

	krokHandler := krok.NewHookHandler(krok.Config{}, krok.Dependencies{
		Logger: log,
	})

	// Create server
	server := server.NewKrokServer(krokArgs.server, server.Dependencies{
		Logger: log,
		Krok:   krokHandler,
	})

	// Run service & server
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return server.Run(ctx) })
	if err := g.Wait(); err != nil {
		log.Err(err).Msg("Failed to run")
	}
}
