// Package command contains cli commands.
package command

import (
	gzGlue "github.com/gozix/glue/v2"
	gzViper "github.com/gozix/viper/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/pwd"
	"log"
	"os"
	"sync"
)

type loader struct {
	logger  *log.Logger
	config  *gzViper.Viper
	dir     string
	loaders []domain.Loader
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
			return &cobra.Command{
				Use:           "load",
				Short:         "Loads data to neo4j",
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
			config = ctn.Get(gzViper.BundleName).(*viper.Viper)
			dir    = ctn.Get(pwd.BundleName).(*pwd.PWD).CurrentDir()
			logger = log.New(os.Stdout, "", log.LstdFlags)
		)
		config = config.Sub("app.imdb")
		if config == nil {
			logger.Println("app.imdb not exist")
			return ErrKeyNotExist
		}

		var d = loader{
			logger:  logger,
			config:  config,
			dir:     dir,
			loaders: make([]domain.Loader, 0),
		}

		for _, def := range ctn.Definitions() {
			for _, tag := range def.Tags {
				if tag.Name == domain.DefTagLoader {
					d.loaders = append(d.loaders, ctn.Get(def.Name).(domain.Loader))
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
		go func(loader domain.Loader) {
			defer wg.Done()

			if err := loader.Load(cmd.Context()); err != nil {
				s.logger.Printf("%s loader err: %s", loader.Name(), err)
			}
		}(s.loaders[i])
	}
	wg.Wait()
	s.logger.Println("Finished!")

	return
}
