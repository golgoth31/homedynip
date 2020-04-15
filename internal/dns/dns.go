// Package dns is used to dyndns provider
package dns

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/spf13/viper"

	"github.com/golgoth31/homedynip/internal/dns/noip"
	"github.com/golgoth31/homedynip/internal/dns/ovh"
)

// New returns a new DNS provider object
func New(conf *viper.Viper, ip string, log *zerolog.Logger) (DNS, error) {
	switch conf.GetString("client.dns") {
	case "ovh":
		return &ovh.Ovh{
			Username: conf.GetString("ovh.username"),
			Password: conf.GetString("ovh.password"),
			Hostname: conf.GetString("ovh.hostname"),
			IP:       ip,
			Log:      log.With().Str("dns-provider", "ovh").Logger(),
		}, nil
	case "noip":
		return &noip.Noip{
			Username: conf.GetString("noip.username"),
			Password: conf.GetString("noip.password"),
			Hostname: conf.GetString("noip.hostname"),
			IP:       ip,
			Log:      log.With().Str("dns-provider", "noip").Logger(),
		}, nil
	default:
		return &dummy{}, fmt.Errorf("unknow DNS driver: %s", conf.GetString("client.dns"))
	}
}
