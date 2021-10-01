// Package command contains cli commands.
package command

import (
	gzGlue "github.com/gozix/glue/v2"
	gzViper "github.com/gozix/viper/v2"
	gzZap "github.com/gozix/zap/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// DefCommandLoadName is container name.
const DefCommandLoadName = "cli.command.message"

// DefCommandLoad register command in di container.
func DefCommandLoad() di.Def {
	return di.Def{
		Name: DefCommandLoadName,
		Tags: []di.Tag{{
			Name: gzGlue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &cobra.Command{
				Use:           "message",
				Short:         "Write configured message to log",
				SilenceUsage:  true,
				SilenceErrors: true,
				RunE: func(cmd *cobra.Command, args []string) error {

					var cfg *viper.Viper
					if err = ctn.Fill(gzViper.BundleName, &cfg); err != nil {
						return err
					}

					var logger *zap.Logger
					if err = ctn.Fill(gzZap.BundleName, &logger); err != nil {
						return err
					}

					return nil
				},
			}, nil
		},
	}
}
