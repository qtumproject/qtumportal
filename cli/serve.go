package cli

import (
	"github.com/hayeah/qtum-portal"

	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	cmd := app.Command("serve", "Start DApp server")

	// bind := cmd.Flag("bind", "net interface to bind server ports").Default("127.0.0.1").String()
	dappPort := cmd.Flag("dapp-port", "port to serve DApp content").Default("9888").Int()
	authPort := cmd.Flag("auth-port", "port to serve DApp authorization API").Default("9899").Int()

	appDir := cmd.Flag("appdir", "DApp content directory").Default(".").String()

	cmd.Action(func(pc *kingpin.ParseContext) error {
		opts := portal.ServerOption{
			DAppPort:      *dappPort,
			AuthPort:      *authPort,
			StaticBaseDir: *appDir,

			QtumdRPCURL: getQtumRPCURL(),
		}

		s := portal.NewServer(opts)

		return s.Start()
	})
}
