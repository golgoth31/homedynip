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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const badExitCode = 1

var logger zerolog.Logger
var config *viper.Viper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "homedynip",
	Short: "Homedynip",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(badExitCode)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	rootCmd.PersistentFlags().String("logLevel", "error", "config file")

	if err := viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("logLevel")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag logLevel: %v", err)
	}

	rootCmd.PersistentFlags().String("logFormat", "console", "config file")

	if err := viper.BindPFlag("logFormat", rootCmd.PersistentFlags().Lookup("logFormat")); err != nil {
		logger.Fatal().Err(err).Msgf("Can't bind flag logFormat: %v", err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() { //nolint:gocyclo
	replacer := strings.NewReplacer(".", "_")
	config = viper.GetViper()

	config.SetEnvPrefix("HOMEDYNIP")
	config.SetEnvKeyReplacer(replacer)
	config.AutomaticEnv()

	if cfgFile != "" {
		config.SetConfigFile(cfgFile)

		err := config.ReadInConfig()
		if err != nil {
			logger.Fatal().Err(err).Msgf("Can't read config: %v", err)
		}
	}

	switch config.GetString("logLevel") {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	switch config.GetString("logFormat") {
	case "json":
		logger = log.With().Logger()
	default:
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}
}
