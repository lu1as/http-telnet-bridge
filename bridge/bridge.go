package bridge

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Cristofori/kmud/telnet"
)

const tcpTimeout = time.Second

type Bridge struct {
	telnet     *telnet.Telnet
	authSecret string
}

func NewBridge(tcpAddr string) (*Bridge, error) {
	c, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	log.Printf("connected to tcp address %s", c.RemoteAddr())

	return &Bridge{
		telnet: telnet.NewTelnet(c),
	}, nil
}

func (b *Bridge) Start(httpAddr string, authSecret string, cert string, key string) error {
	b.authSecret = authSecret
	http.HandleFunc("/ping", b.handlePing)
	http.HandleFunc("/json", b.handleJSON)

	log.Printf("listen on %s for http requests", httpAddr)
	return http.ListenAndServeTLS(httpAddr, cert, key, nil)
}

func (b *Bridge) Stop() {
	b.telnet.Close()
}

func (b *Bridge) handleJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, JsonError("method not implemented"), 400)
		return
	} else if t, ok := r.Header["Authorization"]; !ok || t[0] != b.authSecret {
		http.Error(w, JsonError("unauthorized"), 400)
		return
	} else if t, ok := r.Header["Content-Type"]; !ok || t[0] != "application/json" {
		http.Error(w, JsonError("invalid content type"), 500)
		return
	}

	pl, _ := ioutil.ReadAll(r.Body)
	res, n, err := b.forward(pl)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, JsonError("forwarding failed"), 500)
		return
	}

	w.Write(res[:n])
}

func (b *Bridge) handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I am alive!")
}

func (b *Bridge) forward(payload []byte) (res []byte, n int, err error) {
	b.telnet.SetWriteDeadline(time.Now().Add(tcpTimeout))
	if _, err := b.telnet.Write(append(payload, 0x22, 0x0a)); err != nil {
		return nil, 0, err
	}

	res = make([]byte, 4096)
	b.telnet.SetReadDeadline(time.Now().Add(tcpTimeout))
	n, err = b.telnet.Read(res)
	if err != nil {
		return nil, 0, err
	}

	return res, n, nil
}
