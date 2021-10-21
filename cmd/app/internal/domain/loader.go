package domain

import "context"

const DefTagLoader = "loader"

type Loader interface {
	Name() string
	Load(ctx context.Context) error
}
