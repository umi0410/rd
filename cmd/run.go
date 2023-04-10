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
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"rd/metric"
	"rd/repository"
	"rd/repository/store"
	"rd/server"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run rd server",
	Long: `run rd server. You are recommended to add the server to the 
search engine list of your web browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	host := runCmd.Flags().StringP("host", "h", "0.0.0.0", "Host to listen")
	port := runCmd.Flags().StringP("port", "p", "18081", "Port to listen")

	st, err := store.NewLocalStore()
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}

	aliasRepository := &repository.LocalAliasRepository{
		Store: st,
	}

	s, err := server.NewServer(aliasRepository)
	if err != nil {
		log.Errorf("%+v", err)
	}

	go metric.Run()

	if err := s.Run(*host, *port); err != nil {
		log.Errorf("%+v", err)
	}
}
