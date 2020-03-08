package build

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugcmd.Aliaser = &Cmd{}
var _ plugcmd.SubCommander = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}

type Cmd struct {
	// Mod is the -mod flag
	mod string
	// Static sets the following flags for the final `go build` command:
	// -linkmode external
	// -extldflags "-static"
	static bool
	// Environment the binary is meant for. defaults to "development"
	environment string
	// LDFlags to be passed to the final `go build` command
	ldFlags string
	// BuildFlags to be passed to the final `go build` command
	buildFlags             []string
	tags                   string
	bin                    string
	verbose                bool
	skipTemplateValidation bool

	help      bool
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (*Cmd) CmdAliases() []string {
	return []string{"b", "install"}
}

func (b Cmd) PluginName() string {
	return "build"
}

func (Cmd) Description() string {
	return "Build the application binary, including bundling of webpack (packr & webpack)"
}

func (bc *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range bc.ScopedPlugins() {
		if _, ok := p.(Builder); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (bc *Cmd) ScopedPlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	if bc.pluginsFn != nil {
		plugs = bc.pluginsFn()
	}

	var builders []plugins.Plugin
	for _, p := range plugs {
		switch p.(type) {
		case BuildArger:
			builders = append(builders, p)
		case Builder:
			builders = append(builders, p)
		case BeforeBuilder:
			builders = append(builders, p)
		case AfterBuilder:
			builders = append(builders, p)
		case Versioner:
			builders = append(builders, p)
		case TemplatesValidator:
			builders = append(builders, p)
		case Packager:
			builders = append(builders, p)
		case PackFiler:
			builders = append(builders, p)
		case Flagger:
			builders = append(builders, p)
		case Pflagger:
			builders = append(builders, p)
		case Importer:
			builders = append(builders, p)
		case GoBuilder:
			builders = append(builders, p)
		case Stdouter:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stderrer:
			plugs = append(plugs, p)
		}
	}
	return builders
}
