package cmd

import (
	"flag"

	"github.com/krok-o/krok/pkg/server"
)

var (
	rootArgs struct {
		devMode bool
		debug   bool
		server  server.Config
	}
)

func init() {
	flag.BoolVar(&rootArgs.server.AutoTLS, "auto-tls", false, "--auto-tls")
	flag.BoolVar(&rootArgs.debug, "debug", false, "--debug")
	flag.StringVar(&rootArgs.server.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&rootArgs.server.ServerKeyPath, "server-key-path", "", "--server-key-path /home/user/.server/server.key")
	flag.StringVar(&rootArgs.server.ServerCrtPath, "server-crt-path", "", "--server-crt-path /home/user/.server/server.crt")
	flag.StringVar(&rootArgs.server.Port, "port", "9998", "--port 443")
	flag.Parse()
}
