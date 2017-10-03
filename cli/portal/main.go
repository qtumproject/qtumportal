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
		QtumdRPCURL: qtumRPCURL,
	}

	s := portal.Server{
		Options: opts,
	}

	err = s.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
