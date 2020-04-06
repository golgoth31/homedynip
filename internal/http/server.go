package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func NewServer(port int32) *Server {
	return &Server{
		Config: &http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *Server) Start() {
	s.Config.Handler = http.HandlerFunc(s.echoIp)
	log.Fatal(s.Config.ListenAndServe())
}

func (s *Server) echoIp(w http.ResponseWriter, r *http.Request) {

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
		Ip: ip,
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
