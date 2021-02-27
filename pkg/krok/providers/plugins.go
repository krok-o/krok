package providers

import (
	"context"
)

// Plugins defines functionality a plugin provider needs to have.
// @WIP
type Plugins interface {
	// Upload enables the uploading of a plugin which will define a command.
	Upload(ctx context.Context, id string) error
	// Load ??? given a URL load download and zip it.. or maybe have a github link to a code and we build the plugin?
	Load(ctx context.Context) error
}
