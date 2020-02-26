package newapp_test

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
)

var _ cli.NonAppNeeder = newapp.Cmd{}
