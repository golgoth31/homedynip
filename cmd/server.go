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
	homedynip "github.com/golgoth31/homedynip/internal/http"
	"github.com/spf13/cobra"
)

var port int32

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Simple server to send back public IP to client",
	Run: func(cmd *cobra.Command, args []string) {
		serv := homedynip.NewServer(port)
		serv.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().Int32Var(&port, "port", 8080, "port to listen on")
}
