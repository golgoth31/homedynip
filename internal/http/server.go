package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

const timeout = 10

// NewServer returns a new homedynip server
func NewServer() *Server {
	return &Server{
		HTTPConfig: &http.Server{
			ReadTimeout:  timeout * time.Second,
			WriteTimeout: timeout * time.Second,
		},
	}
}

// Start starts the server
func (s *Server) Start() {
	s.HTTPConfig.Addr = fmt.Sprintf(":%d", s.Config.GetInt32("server.port"))
	s.HTTPConfig.Handler = http.HandlerFunc(s.echoIP)

	err := s.HTTPConfig.ListenAndServe()
	if err != nil {
		s.Log.Fatal().Err(err).Msg("")
	}
}

func (s *Server) echoIP(w http.ResponseWriter, r *http.Request) {
	var err error

	var ip string

	if r.Header.Get("X-Forwarded-For") != "" {
		ip = r.Header.Get("X-Forwarded-For")
	} else {
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			s.Log.Error().Err(err).Msgf("userip: %q is not IP:port", r.RemoteAddr)
			ip = "unknown"
		}
	}

	s.Log.Info().Msgf("Request from %s", ip)

	output := response{
		IP: ip,
	}

	jData, err := json.Marshal(output)
	if err != nil {
		s.Log.Error().Err(err).Msgf("%v", err)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jData)
	if err != nil {
		s.Log.Error().Err(err).Msgf("Unable to write response: %v", err)
	}
}
