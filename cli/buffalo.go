package cli

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/clifix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/bzr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/flect"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/git"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/i18n"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/mail"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/packr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/plush"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/refresh"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ plugcmd.SubCommander = &Buffalo{}
var _ plugins.Plugin = &Buffalo{}
var _ plugins.Scoper = &Buffalo{}
var _ plugprint.Describer = &Buffalo{}

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	plugins.Plugins
}

func insidePlugins(root string) []plugins.Plugin {
	var plugs []plugins.Plugin

	plugs = append(plugs, cmds.InsidePlugins()...)
	plugs = append(plugs, clifix.Plugins()...)
	plugs = append(plugs, fizz.Plugins()...)
	plugs = append(plugs, flect.Plugins()...)
	plugs = append(plugs, golang.Plugins()...)
	plugs = append(plugs, grifts.Plugins()...)
	plugs = append(plugs, i18n.Plugins()...)
	plugs = append(plugs, mail.Plugins()...)
	plugs = append(plugs, packr.Plugins()...)
	plugs = append(plugs, pkger.Plugins()...)
	plugs = append(plugs, plush.Plugins()...)
	plugs = append(plugs, pop.Plugins()...)
	plugs = append(plugs, refresh.Plugins()...)
	plugs = append(plugs, soda.Plugins()...)

	if _, err := os.Stat(filepath.Join(root, "package.json")); err == nil {
		plugs = append(plugs, webpack.Plugins()...)
	}

	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		plugs = append(plugs, git.Plugins()...)
	}

	if _, err := os.Stat(filepath.Join(root, ".bzr")); err == nil {
		plugs = append(plugs, bzr.Plugins()...)
	}
	return plugs
}

func outsidePlugins(root string) []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, cmds.OutsidePlugins()...)
	return plugs
}

func NewFromRoot(root string) (*Buffalo, error) {
	b := &Buffalo{}

	isBuffalo := IsBuffalo(root)

	pfn := func() []plugins.Plugin {
		return b.Plugins
	}

	if isBuffalo {
		b.Plugins = append(b.Plugins, insidePlugins(root)...)
	} else {
		b.Plugins = append(b.Plugins, outsidePlugins(root)...)
	}

	plugs := make([]plugins.Plugin, 0, len(b.Plugins))
	for _, p := range b.Plugins {
		switch t := p.(type) {
		case NonAppNeeder:
			if isBuffalo {
				continue
			}
		case AppNeeder:
			if !isBuffalo {
				continue
			}
		case AvailabilityChecker:
			if !t.PluginAvailable(root) {
				continue
			}
		}
		plugs = append(plugs, p)
	}

	b.Plugins = plugs

	sort.Slice(b.Plugins, func(i, j int) bool {
		return b.Plugins[i].PluginName() < b.Plugins[j].PluginName()
	})

	pfn = func() []plugins.Plugin {
		return b.Plugins
	}

	for _, b := range b.Plugins {
		f, ok := b.(plugins.Needer)
		if !ok {
			continue
		}
		f.WithPlugins(pfn)
	}

	return b, nil
}

func New() (*Buffalo, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return NewFromRoot(pwd)
}

func (b Buffalo) ScopedPlugins() []plugins.Plugin {
	return b.Plugins
}

func (b Buffalo) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range b.ScopedPlugins() {
		if _, ok := p.(Commander); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

// Name ...
func (Buffalo) PluginName() string {
	return "buffalo"
}

func (Buffalo) String() string {
	return "buffalo"
}

// Description ...
func (Buffalo) Description() string {
	return "Tools for working with Buffalo applications"
}
