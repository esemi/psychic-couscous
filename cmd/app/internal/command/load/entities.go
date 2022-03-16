// Package load contains cli commands.
package load

import (
	"log"
	"os"
	"sync"

	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/pwd"
)

type loadEntities struct {
	logger  *log.Logger
	dir     string
	loaders []domain.LoaderEntities
}

// DefCommandLoadEntitiesName is container name.
const DefCommandLoadEntitiesName = "cli.command.load-entities"

// DefCommandLoadEntities register command in di container.
func DefCommandLoadEntities() di.Def {
	return di.Def{
		Name: DefCommandLoadEntitiesName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &cobra.Command{
				Use:           "entities",
				Short:         "load entities data to neo4j",
				SilenceUsage:  true,
				SilenceErrors: true,
				RunE:          EntitiesRunE(ctn),
			}, nil
		},
	}
}

func EntitiesRunE(ctn di.Container) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		var (
			dir    = ctn.Get(pwd.BundleName).(*pwd.PWD).CurrentDir()
			logger = log.New(os.Stdout, "", log.LstdFlags)
		)

		var d = loadEntities{
			logger:  logger,
			dir:     dir,
			loaders: make([]domain.LoaderEntities, 0),
		}

		for _, def := range ctn.Definitions() {
			for _, tag := range def.Tags {
				if tag.Name == domain.DefTagLoaderEntities {
					d.loaders = append(d.loaders, ctn.Get(def.Name).(domain.LoaderEntities))
				}
			}
		}

		return d.Handler(cmd, args)
	}
}

// Handler run.
func (s *loadEntities) Handler(cmd *cobra.Command, _ []string) (err error) {
	s.logger.Println("Start...")
	var wg sync.WaitGroup
	for i := range s.loaders {
		wg.Add(1)
		go func(loader domain.LoaderEntities) {
			defer wg.Done()

			if err := loader.LoadEntities(cmd.Context()); err != nil {
				s.logger.Printf("%s load entities err: %s", loader.Name(), err)
			}
		}(s.loaders[i])
	}
	wg.Wait()
	s.logger.Println("Finished!")

	return
}
