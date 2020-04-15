package ovh

import "github.com/rs/zerolog"

// Ovh represents an ovh client
type Ovh struct {
	Username string
	Password string
	Hostname string
	IP       string
	Log      zerolog.Logger
}
