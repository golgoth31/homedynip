// Package cmd ...
/*
Copyright © 2020 David Sabatie <david.sabatie@notrenet.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"time"

	homedynip "github.com/golgoth31/homedynip/internal/http"
	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Request my own IP address and send it a dyndns server",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			if config.GetBool("client.showip") {
				zerolog.SetGlobalLevel(zerolog.FatalLevel)
			}
			client := homedynip.NewClient()
			client.Log = &logger
			client.Config = config

			output, err := client.GetIP()
			if err != nil {
				logger.Error().Err(err).Msgf("Can't get IP: %v", err)
			}
			logger.Info().Msgf("My Ip is: %s", output)
			if !config.GetBool("client.dryrun") {
				if err := client.WriteDNS(); err != nil {
					logger.Error().Err(err).Msgf("Unable to write IP: %v", err)
				}
			}
			if config.GetBool("client.cron") {
				break
			}
			logger.Info().Msgf("sleeping for %v", config.GetDuration("client.sleep"))
			time.Sleep(config.GetDuration("client.sleep"))
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().String("url", "", "URL of custom service")

	if err := viper.BindPFlag("client.url", clientCmd.PersistentFlags().Lookup("url")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag url: %v", err)
	}

	clientCmd.PersistentFlags().String("service", "", "Service to request")

	if err := viper.BindPFlag("client.service", clientCmd.PersistentFlags().Lookup("service")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag url: %v", err)
	}

	clientCmd.PersistentFlags().Bool("insecure", false, "Insecure https request")

	if err := viper.BindPFlag("client.insecure", clientCmd.PersistentFlags().Lookup("insecure")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag url: %v", err)
	}

	clientCmd.PersistentFlags().Bool("cron", false, "Run as cronjob (external scheduling)")

	if err := viper.BindPFlag("client.cron", clientCmd.PersistentFlags().Lookup("cron")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag url: %v", err)
	}

	duration, err := time.ParseDuration("24h")
	if err != nil {
		logger.Fatal().Err(err).Msgf("Unable to parse duration: %v", err)
	}

	clientCmd.PersistentFlags().Duration("sleep", duration, "Wait for duration between 2 IP requests")

	if err := viper.BindPFlag("client.sleep", clientCmd.PersistentFlags().Lookup("sleep")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag url: %v", err)
	}

	clientCmd.PersistentFlags().Bool("dryrun", false, "Only request IP service")

	if err := viper.BindPFlag("client.dryrun", clientCmd.PersistentFlags().Lookup("dryrun")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag dryrun: %v", err)
	}

	clientCmd.PersistentFlags().Bool("showip", false, "Only print out Received IP")

	if err := viper.BindPFlag("client.showip", clientCmd.PersistentFlags().Lookup("showip")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag showip: %v", err)
	}
}
