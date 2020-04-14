package noip

import "github.com/rs/zerolog"

// Noip represents a noip client
type Noip struct {
	Username string
	Password string
	Hostname string
	IP       string
	Log      *zerolog.Logger
}
