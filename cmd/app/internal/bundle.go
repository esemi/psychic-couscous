// Package internal provide component implementation.
package internal

import (
	gzViper "github.com/gozix/viper/v2"
	gzZap "github.com/gozix/zap/v2"
	"github.com/sarulabs/di/v2"

	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/command"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/controller"
)

// Bundle is component bundle.
type Bundle struct{}

// NewBundle is bundle constructor.
func NewBundle() *Bundle {
	return &Bundle{}
}

// Name implements the glue.Bundle interface.
func (*Bundle) Name() string {
	return "app"
}

// Build implements the glue.Bundle interface.
func (*Bundle) Build(builder *di.Builder) error {
	return builder.Add(
		// controller
		controller.DefSystemController(),

		// command
		command.DefCommandDownload(),
		command.DefCommandLoad(),
	)
}

// DependsOn implements the glue.BundleDependsOn interface.
func (*Bundle) DependsOn() []string {
	return []string{gzViper.BundleName, gzZap.BundleName}
}
