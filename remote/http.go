package remote

import (
	"fmt"
	"log"
	"net/http"
	"sanHeRecruitment/config"
)

const (
	defaultBasePath = config.DefaultBasePath
)

var defaultHttpIp = config.DefaultHttpIp

var _ ToService = (*HttpToService)(nil)

// Server  HTTPPool implements PeerPicker for a pool of HTTP peers.
type Server struct {
	self     string
	basePath string
}

// NewHTTPPool initializes an HTTP pool of peers
func newHTTPServer() *Server {
	return &Server{
		self:     defaultHttpIp,
		basePath: defaultBasePath,
	}
}

func StartBeToServer(addr string) {
	HP := newHTTPServer()
	log.Println("sanheServer is running at", addr)
	log.Fatal(http.ListenAndServe(addr, HP))
}

// Log info with server name
func (p *Server) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s \n", p.self, fmt.Sprintf(format, v...))
}
