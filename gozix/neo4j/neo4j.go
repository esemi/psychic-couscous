package neo4j

import (
	"fmt"

	gzviper "github.com/gozix/viper/v2"
	neo "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sarulabs/di/v2"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	Driver = neo.Driver
)

// BundleName is default definition name.
const BundleName = "neo4j.driver"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Name implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	return builder.Add(
		di.Def{
			Name: BundleName,
			Build: func(ctn di.Container) (_ interface{}, err error) {
				var config *gzviper.Viper
				if err = ctn.Fill(gzviper.BundleName, &config); err != nil {
					return nil, err
				}

				var suffix = fmt.Sprintf("%s.", "neo4j")

				var (
					addr  = config.GetString(suffix + "addr")
					user  = config.GetString(suffix + "username")
					pass  = config.GetString(suffix + "password")
					realm = config.GetString(suffix + "realm")
				)
				dbUri := fmt.Sprintf("neo4j://%s", addr)
				driver, err := neo.NewDriver(dbUri, neo.BasicAuth(user, pass, realm))
				if err != nil {
					return nil, err
				}

				return driver, nil
			},
			Close: func(obj interface{}) error {
				return obj.(neo.Driver).Close()
			},
		},
	)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{gzviper.BundleName}
}
