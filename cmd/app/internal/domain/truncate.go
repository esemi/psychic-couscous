package domain

import "context"

const (
	DefTagTruncateEntities = "truncate-entities"
)

type TruncaterEntities interface {
	Name() string
	Truncate(ctx context.Context) error
}
