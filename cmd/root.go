package cmd

import (
	"github.com/spf13/cobra"

	"github.com/krok-o/krok/pkg/krok/providers/environment"
	"github.com/krok-o/krok/pkg/krok/providers/livestore"
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

}

func runKrokCmd(cmd *cobra.Command, args []string) {
}
