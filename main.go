package main

import (
	"log"

	"github.com/lu1as/http-telnet-bridge/bridge"
	"github.com/namsral/flag"
)

var (
	httpAddr   string
	httpCert   string
	httpKey    string
	authSecret string
	tcpAddr    string
)

func init() {
	flag.StringVar(&httpAddr, "httpAddr", ":8080", "ip and port used for serving http")
	flag.StringVar(&httpCert, "httpCert", "server.crt", "ssl certificate")
	flag.StringVar(&httpKey, "httpKey", "server.key", "ssl certificate key")
	flag.StringVar(&tcpAddr, "tcpAddr", ":1705", "ip and port of tcp endpoint forwarding to")
	flag.StringVar(&authSecret, "authSecret", "secret", "authentication secret for requests")
}

func main() {
	flag.Parse()

	b, err := bridge.NewBridge(tcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer b.Stop()
	if err := b.Start(httpAddr, authSecret, httpCert, httpKey); err != nil {
		log.Fatal(err.Error())
	}
}
