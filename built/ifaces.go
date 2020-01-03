package built

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
)

// Initer is invoked in when an application binary
// built with `buffalo build` is executed. This hook
// is executed before any flags are parsed or sub-commands
// are run.
type Initer interface {
	plugins.Plugin
	BuiltInit(ctx context.Context, args []string) error
}