package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const timeout = 10

// NewServer returns a new homedynip server
func NewServer(port int32) *Server {
	return &Server{
		Config: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  timeout * time.Second,
			WriteTimeout: timeout * time.Second,
		},
	}
}

// Start starts the server
func (s *Server) Start() {
	s.Config.Handler = http.HandlerFunc(s.echoIP)
	log.Fatal(s.Config.ListenAndServe())
}

func (s *Server) echoIP(w http.ResponseWriter, r *http.Request) {
	var err error

	var ip string

	if r.Header.Get("X-Forwarded-For") != "" {
		ip = r.Header.Get("X-Forwarded-For")
	} else {
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("userip: %q is not IP:port", r.RemoteAddr)
			ip = "unknown"
		}
	}

	log.Printf("Request from %s", ip)

	output := response{
		IP: ip,
	}

	jData, err := json.Marshal(output)
	if err != nil {
		log.Printf("%v", err)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jData)
	if err != nil {
		log.Printf("Unable to write response: %v", err)
	}
}
