// Package command contains cli commands.
package command

import (
	"fmt"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"

	gzGlue "github.com/gozix/glue/v2"
	gzViper "github.com/gozix/viper/v2"
	gzZap "github.com/gozix/zap/v2"
)

// DefCommandMessageName is container name.
const DefCommandMessageName = "cli.command.message"

// DefCommandMessage register command in di container.
func DefCommandMessage() di.Def {
	return di.Def{
		Name: DefCommandMessageName,
		Tags: []di.Tag{{
			Name: gzGlue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &cobra.Command{
				Use:           "message",
				Short:         "Write configured message to log",
				SilenceUsage:  true,
				SilenceErrors: true,
				Args:          cobra.ExactArgs(1),
				RunE: func(cmd *cobra.Command, args []string) error {
					fmt.Print(os.Args)

					var cfg *viper.Viper
					if err = ctn.Fill(gzViper.BundleName, &cfg); err != nil {
						return err
					}

					var logger *zap.Logger
					if err = ctn.Fill(gzZap.BundleName, &logger); err != nil {
						return err
					}

					logger.Info(cfg.GetString("message"))

					return nil
				},
			}, nil
		},
	}
}
