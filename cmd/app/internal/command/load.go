// Package command contains cli commands.
package command

import (
	"log"
	"os"

	gzGlue "github.com/gozix/glue/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/command/load"
)

type loader struct {
	logger *log.Logger
	cmds   []*cobra.Command
}

// DefCommandLoadName is container name.
const DefCommandLoadName = "cli.command.load"

// DefCommandLoad register command in di container.
func DefCommandLoad() di.Def {
	return di.Def{
		Name: DefCommandLoadName,
		Tags: []di.Tag{{
			Name: gzGlue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cmd = &cobra.Command{
				Use:           "load",
				Short:         "Loads data to neo4j",
				SilenceUsage:  true,
				SilenceErrors: true,
				RunE:          LoadRunE(ctn),
			}

			cmd.AddCommand(ctn.Get(load.DefCommandLoadEntitiesName).(*cobra.Command))
			cmd.AddCommand(ctn.Get(load.DefCommandLoadRelationsName).(*cobra.Command))

			return cmd, nil
		},
	}
}

func LoadRunE(ctn di.Container) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		var logger = log.New(os.Stdout, "", log.LstdFlags)

		var d = loader{
			logger: logger,
			cmds: []*cobra.Command{
				ctn.Get(load.DefCommandLoadEntitiesName).(*cobra.Command),
				ctn.Get(load.DefCommandLoadRelationsName).(*cobra.Command),
			},
		}

		return d.Handler(cmd, args)
	}
}

// Handler run.
func (s *loader) Handler(cmd *cobra.Command, args []string) (err error) {
	for _, c := range s.cmds {
		s.logger.Printf("Command: %s \n Err: %s", c.Name(), c.RunE(cmd, args))
	}

	return
}
