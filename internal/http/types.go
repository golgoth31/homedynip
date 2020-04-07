package http

import (
	"net"
	"net/http"

	"github.com/rs/zerolog"
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
	Log    *zerolog.Logger
}

// Server represent an homedynip server
type Server struct {
	HTTPConfig *http.Server
	Config     *viper.Viper
	Log        *zerolog.Logger
}

type response struct {
	IP string `json:"ip"`
}
