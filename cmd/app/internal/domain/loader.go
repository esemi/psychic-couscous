package domain

import "context"

const (
	DefTagLoaderRelations = "loader-relations"
	DefTagLoaderEntities  = "loader-entities"
)

type LoaderRelations interface {
	Name() string
	LoadRelations(ctx context.Context) error
}

type LoaderEntities interface {
	Name() string
	LoadEntities(ctx context.Context) error
}
