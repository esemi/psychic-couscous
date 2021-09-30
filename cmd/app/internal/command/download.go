// Package command contains cli commands.
package command

import (
	gzGlue "github.com/gozix/glue/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
)

// DefCommandDownloadName is container name.
const DefCommandDownloadName = "cli.command.cookie"

// DefCommandDownload register command in di container.
func DefCommandDownload() di.Def {
	return di.Def{
		Name: DefCommandDownloadName,
		Tags: []di.Tag{{
			Name: gzGlue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cmd = &cobra.Command{
				Use:           "download",
				Short:         "download imdb data",
				SilenceUsage:  true,
				SilenceErrors: true,
				Args:          cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) (err error) {

					return nil
				},
			}

			return cmd, nil
		},
	}
}
