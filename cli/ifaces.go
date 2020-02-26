package cli

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
)

type Aliaser = plugcmd.Aliaser
type Commander = plugcmd.Commander
type Needer = plugins.Needer
type StderrNeeder = plugio.ErrNeeder
type StdinNeeder = plugio.InNeeder
type StdoutNeeder = plugio.OutNeeder

type AvailabilityChecker interface {
	PluginAvailable(root string) bool
}

// AppNeeder plugins will only be available inside
// applications that require buffalo.
type AppNeeder interface {
	plugins.Plugin
	NeedsBuffaloApp()
}

// NonAppNeeder plugins will be available
// outside of buffalo applications
type NonAppNeeder interface {
	plugins.Plugin
	NoBuffaloApp()
}
