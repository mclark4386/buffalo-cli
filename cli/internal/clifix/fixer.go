package clifix

import (
	"context"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/fix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/cligen"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ plugins.Plugin = &Fixer{}
var _ plugcmd.Namer = &Fixer{}
var _ fix.Fixer = &Fixer{}

type Fixer struct {
}

func (*Fixer) PluginName() string {
	return "cli/fixer"
}

func (*Fixer) CmdName() string {
	return "cli"
}

func (fixer *Fixer) Fix(ctx context.Context, root string, args []string) error {
	g := &cligen.Generator{}
	return g.Generate(ctx, root, args)
}
