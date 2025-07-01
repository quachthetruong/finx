package repository

import "context"

type FlexOpenApiRepository interface {
	IsHOActive(ctx context.Context) (bool, error)
}
