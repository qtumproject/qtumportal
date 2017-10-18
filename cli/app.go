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
	if *qtumRPC == "" {
		log.Fatalln("Please set QTUM_RPC to qtumd's RPC URL")
	}

	url, err := url.Parse(*qtumRPC)
	if err != nil {
		log.Fatalln("QTUM_RPC URL:", *qtumRPC)
	}

	if url.User == nil {
		log.Fatalln("QTUM_RPC URL (must specify user & password):", *qtumRPC)
	}

	return url
}
