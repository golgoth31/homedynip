package http

import (
	"net"
	"net/http"

	"github.com/spf13/viper"
)

type Client struct {
	// Service  string
	// URL      string
	// Insecure bool
	// Cron     bool
	// Sleep    time.Duration
	Config *viper.Viper
	Ip     *net.IPAddr
}

type Server struct {
	Config *http.Server
}

type response struct {
	Ip string `json:"ip"`
}
