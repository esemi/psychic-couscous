// Package main provide component entry point.
package main

import (
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/neo4j"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/pwd"
	"log"

	gzEcho "github.com/gozix/echo/v2"
	gzGlue "github.com/gozix/glue/v2"
	gzUT "github.com/gozix/universal-translator/v2"
	gzValidator "github.com/gozix/validator/v2"
	gzViper "github.com/gozix/viper/v2"
	gzZap "github.com/gozix/zap/v2"

	gzInternal "gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal"
)

// Version is component version.
const Version = "0.0.1"

func main() {
	var app, err = gzGlue.NewApp(
		gzGlue.Version(Version),
		gzGlue.Bundles(
			gzValidator.NewBundle(),
			gzViper.NewBundle(),
			gzZap.NewBundle(),
			gzInternal.NewBundle(),
			gzEcho.NewBundle(),
			gzUT.NewBundle(),
			pwd.NewBundle(),
			neo4j.NewBundle(),
		),
	)

	if err != nil {
		log.Fatalf("Some error occurred during create app. Error: %v\n", err)
	}

	if err = app.Execute(); err != nil {
		log.Fatalf("Some error occurred during execute app. Error: %v\n", err)
	}
}
