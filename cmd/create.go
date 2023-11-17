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
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"rd/entity"
)

var (
	group       string
	name        string
	destination string
)

// createCmd represents the create command
var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			_, _, aliasService, _ := initialize()
			alias, err := aliasService.Create(context.Background(), &entity.Alias{
				AliasGroup:  group,
				Name:        name,
				Destination: destination,
			})
			if err != nil {
				log.Panicf("%+v", err)
			}
			log.Infof("alias \"%s/%s\" has been create.", alias.AliasGroup, alias.Name)
		},
	}
)

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	createCmd.PersistentFlags().StringVarP(&group, "group", "g", "", "The group of the alias to create")
	createCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "The name of the alias to create")
	createCmd.Flags().StringVarP(&destination, "dst", "d", "", "The destination of the alias to create")
	if err := createCmd.MarkPersistentFlagRequired("group"); err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	if err := createCmd.MarkPersistentFlagRequired("name"); err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	if err := createCmd.MarkFlagRequired("dst"); err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
}
