// Package command contains cli commands.
package command

import (
	gzGlue "github.com/gozix/glue/v2"
	gzZap "github.com/gozix/zap/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/database"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"go.uber.org/zap"
)

// DefCommandTruncateName is container name.
const DefCommandTruncateName = "cli.command.truncate"

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
				RunE: func(cmd *cobra.Command, args []string) (err error) {

					var movieRepository domain.MovieRepository
					if err = ctn.Fill(database.DefMoviesRepositoryName, &movieRepository); err != nil {
						return err
					}

					var logger *zap.Logger
					if err = ctn.Fill(gzZap.BundleName, &logger); err != nil {
						return err
					}

					if err = movieRepository.Truncate(); err != nil {
						return err
					}

					logger.Info("Truncated success")
					return nil
				},
			}, nil
		},
	}
}
