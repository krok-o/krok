package providers

import (
	"context"
)

// ApiKeys defines operations that an api key provider must have.
type ApiKeys interface {
	Generate(ctx context.Context)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context)
	Get(ctx context.Context)
}
