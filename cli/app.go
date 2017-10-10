package cli

import (
	"log"
	"net/url"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var app = kingpin.New("qtumportal", "QTUM DApp Server")

var qtumRPC = app.Flag("qtum-rpc", "URL of qtum RPC service").Envar("QTUM_RPC").Default("").String()

func Run() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func getQtumRPCURL() *url.URL {
	url, err := url.Parse(*qtumRPC)
	if err != nil {
		log.Fatalln("Invalid QTUM RPC URL:", *qtumRPC)
	}
	return url
}
