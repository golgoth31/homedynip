package http

import (
	"net"
	"net/http"

	"github.com/spf13/viper"
)

// Client represent an homedynip client
type Client struct {
	// Service  string
	// URL      string
	// Insecure bool
	// Cron     bool
	// Sleep    time.Duration
	Config *viper.Viper
	IP     *net.IPAddr
}

// Server represent an homedynip server
type Server struct {
	Config *http.Server
}

type response struct {
	IP string `json:"ip"`
}
