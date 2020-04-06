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
	"log"
	"strings"
	"time"

	homedynip "github.com/golgoth31/homedynip/internal/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientConfig *viper.Viper
var cfgFile string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Request my own IP adress and send it a dyndns server",
	Run: func(cmd *cobra.Command, args []string) {
		client := homedynip.NewClient()
		client.Config = clientConfig

		output, err := client.GetIp()
		if err != nil {
			log.Printf("Can't get IP: %v", err)
		}
		log.Printf("My Ip is: %s", output)
		if err := client.WriteDNS(); err != nil {
			log.Printf("Unable to write IP: %v", err)
		}
	},
}

func init() {
	clientConfig = viper.New()

	cobra.OnInitialize(initClientConfig)
	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	clientCmd.PersistentFlags().String("url", "", "URL of custom service")
	if err := clientConfig.BindPFlag("client.url", clientCmd.PersistentFlags().Lookup("url")); err != nil {
		log.Fatalf("Can't bind flag url: %v", err)
	}
	clientCmd.PersistentFlags().String("service", "", "Service to request")
	if err := clientConfig.BindPFlag("client.service", clientCmd.PersistentFlags().Lookup("service")); err != nil {
		log.Fatalf("Can't bind flag url: %v", err)
	}
	clientCmd.PersistentFlags().Bool("insecure", false, "insecure https request")
	if err := clientConfig.BindPFlag("client.insecure", clientCmd.PersistentFlags().Lookup("insecure")); err != nil {
		log.Fatalf("Can't bind flag url: %v", err)
	}
	clientCmd.PersistentFlags().Bool("cron", false, "run as cronjob (external scheduling)")
	if err := clientConfig.BindPFlag("client.cron", clientCmd.PersistentFlags().Lookup("cron")); err != nil {
		log.Fatalf("Can't bind flag url: %v", err)
	}
	duration, _ := time.ParseDuration("24h")
	clientCmd.PersistentFlags().Duration("sleep", duration, "Wait duration between 2 IP requests")
	if err := clientConfig.BindPFlag("client.sleep", clientCmd.PersistentFlags().Lookup("sleep")); err != nil {
		log.Fatalf("Can't bind flag url: %v", err)
	}

}

// initConfig reads in config file and ENV variables if set.
func initClientConfig() {
	replacer := strings.NewReplacer(".", "_")
	clientConfig.SetEnvPrefix("HOMEDYNIP")
	clientConfig.SetEnvKeyReplacer(replacer)
	clientConfig.AutomaticEnv()

	if cfgFile != "" {
		// Use config file from the flag.
		clientConfig.SetConfigFile(cfgFile)

		err := clientConfig.ReadInConfig()
		if err != nil {
			log.Fatalf("Can't read config: %v", err)
		}
	}
}
