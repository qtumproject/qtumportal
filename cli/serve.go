package cli

import (
	"github.com/hayeah/qtum-portal"

	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	cmd := app.Command("serve", "Start DApp server")

	bind := cmd.Flag("bind", "network interface to bind to (e.g. 0.0.0.0) ").Default("").String()

	// bind := cmd.Flag("bind", "net interface to bind server ports").Default("127.0.0.1").String()
	dappPort := cmd.Flag("dapp-port", "port to serve DApp content").Default("9888").Int()
	authPort := cmd.Flag("auth-port", "port to serve DApp authorization API").Default("9899").Int()

	devMode := cmd.Flag("dev", "[Insecure] Developer mode").Default("false").Bool()

	appDir := cmd.Arg("appdir", "DApp content directory").Default("").String()

	cmd.Action(func(pc *kingpin.ParseContext) error {
		staticDir := *appDir

		if !*devMode && staticDir == "" {
			staticDir = "."
		}

		opts := portal.ServerOption{
			Bind:          *bind,
			DAppPort:      *dappPort,
			AuthPort:      *authPort,
			StaticBaseDir: staticDir,

			QtumdRPCURL: getQtumRPCURL(),
			DebugMode:   *devMode,
		}

		s := portal.NewServer(opts)

		return s.Start()
	})
}
