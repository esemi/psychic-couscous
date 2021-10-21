package pwd

import (
	"os"
	"path/filepath"

	"github.com/sarulabs/di/v2"
)

// BundleName is default definition name.
const BundleName = "pwd"

type (
	PWD struct {
		currentDir string
	}
)

func (p *PWD) CurrentDir() string {
	return p.currentDir
}

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}
)

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
				dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					return nil, err
				}

				return &PWD{currentDir: dir}, nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
	)
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{}
}
