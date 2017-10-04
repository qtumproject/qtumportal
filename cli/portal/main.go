package main

import "github.com/hayeah/qtum-portal"
import "log"
import "net/url"
import "os"

func main() {
	qtumRPCURL, err := url.Parse("http://howard:yeh@localhost:13889")
	if err != nil {
		log.Fatalln(err)
	}

	debug, _ := os.LookupEnv("DEBUG")

	opts := portal.ServerOption{
		Port:        9999,
		AuthPort:    9898,
		QtumdRPCURL: qtumRPCURL,

		DebugMode: debug == "true",
	}

	s := portal.NewServer(opts)

	err = s.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
