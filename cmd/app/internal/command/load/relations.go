// Package load contains cli commands.
package load

import (
	"log"
	"os"
	"sync"

	gzViper "github.com/gozix/viper/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/pwd"
)

type loader struct {
	logger  *log.Logger
	config  *gzViper.Viper
	dir     string
	loaders []domain.LoaderRelations
}

// DefCommandLoadRelationsName is container name.
const DefCommandLoadRelationsName = "cli.command.load-relations"

// DefCommandLoadRelations register command in di container.
func DefCommandLoadRelations() di.Def {
	return di.Def{
		Name: DefCommandLoadRelationsName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &cobra.Command{
				Use:           "relations",
				Short:         "Load relations to neo4j",
				SilenceUsage:  true,
				SilenceErrors: true,
				RunE:          LoadRunE(ctn),
			}, nil
		},
	}
}

func LoadRunE(ctn di.Container) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		var (
			dir    = ctn.Get(pwd.BundleName).(*pwd.PWD).CurrentDir()
			logger = log.New(os.Stdout, "", log.LstdFlags)
		)

		var d = loader{
			logger:  logger,
			dir:     dir,
			loaders: make([]domain.LoaderRelations, 0),
		}

		for _, def := range ctn.Definitions() {
			for _, tag := range def.Tags {
				if tag.Name == domain.DefTagLoaderRelations {
					d.loaders = append(d.loaders, ctn.Get(def.Name).(domain.LoaderRelations))
				}
			}
		}

		return d.Handler(cmd, args)
	}
}

// Handler run.
func (s *loader) Handler(cmd *cobra.Command, _ []string) (err error) {
	s.logger.Println("Start...")
	var wg sync.WaitGroup
	for i := range s.loaders {
		wg.Add(1)
		go func(loader domain.LoaderRelations) {
			defer wg.Done()

			if err := loader.LoadRelations(cmd.Context()); err != nil {
				s.logger.Printf("%s loader err: %s", loader.Name(), err)
			}
		}(s.loaders[i])
	}
	wg.Wait()
	s.logger.Println("Finished!")

	return
}
