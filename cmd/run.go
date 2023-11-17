/*
Copyright © 2022 Jinsu Park <dev.umijs@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"rd/metric"
	//"rd/repository/store"
	"rd/server"
)

var (
	host *string
	port *string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run rd server",
	Long: `run rd server. You are recommended to add the server to the 
search engine list of your web browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		aliasRepo, _, aliasService, authService := initialize()

		srv, err := server.NewServer(aliasRepo, aliasService, authService)
		if err != nil {
			log.Errorf("%+v", err)
		}

		go metric.Run()

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		go func() {
			s := <-sigc
			log.Infof("Got SIGNAL, %s", s)
			err := aliasRepo.Close()
			if err != nil {
				log.Panicf("%+v", err)
			}
			os.Exit(0)
		}()
		// TODO: Shutdown webserver gracefully
		if err := srv.Run(*host, *port); err != nil {
			log.Errorf("%+v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// XXX: shorthand h is already defined for help command.
	host = runCmd.Flags().StringP("host", "", "0.0.0.0", "Host to listen")
	port = runCmd.Flags().StringP("port", "p", "18080", "Port to listen")
}
