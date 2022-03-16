// Package command contains cli commands.
package command

import (
	"log"
	"os"
	"sync"

	gzGlue "github.com/gozix/glue/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"

	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
)

// DefCommandTruncateName is container name.
const DefCommandTruncateName = "cli.command.truncate"

type truncateEntities struct {
	logger *log.Logger
	defs   []domain.TruncaterEntities
}

// DefCommandTruncate register command in di container.
func DefCommandTruncate() di.Def {
	return di.Def{
		Name: DefCommandTruncateName,
		Tags: []di.Tag{{
			Name: gzGlue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &cobra.Command{
				Use:           "truncate",
				Short:         "Truncate database",
				SilenceUsage:  true,
				SilenceErrors: true,
				RunE:          TruncateRunE(ctn),
			}, nil
		},
	}
}

func TruncateRunE(ctn di.Container) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		var (
			logger = log.New(os.Stdout, "", log.LstdFlags)
		)

		var d = truncateEntities{
			logger: logger,
			defs:   make([]domain.TruncaterEntities, 0),
		}

		for _, def := range ctn.Definitions() {
			for _, tag := range def.Tags {
				if tag.Name == domain.DefTagTruncateEntities {
					d.defs = append(d.defs, ctn.Get(def.Name).(domain.TruncaterEntities))
				}
			}
		}

		return d.Handler(cmd, args)
	}
}

// Handler run.
func (s *truncateEntities) Handler(cmd *cobra.Command, _ []string) (err error) {
	s.logger.Println("Start...")
	var wg sync.WaitGroup
	for i := range s.defs {
		wg.Add(1)
		go func(loader domain.TruncaterEntities) {
			defer wg.Done()

			if err := loader.Truncate(cmd.Context()); err != nil {
				s.logger.Printf("%s truncate err: %s", loader.Name(), err)
			}
		}(s.defs[i])
	}
	wg.Wait()
	s.logger.Println("Finished!")

	return
}
