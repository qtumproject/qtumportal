package main

import "github.com/hayeah/qtum-portal"
import "log"
import "net/url"

func main() {
	qtumRPCURL, err := url.Parse("http://howard:yeh@localhost:13889")
	if err != nil {
		log.Fatalln(err)
	}

	opts := portal.ServerOption{
		Port:        9999,
		AuthPort:    9898,
		QtumdRPCURL: qtumRPCURL,
	}

	s := portal.NewServer(opts)

	err = s.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
