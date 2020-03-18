package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugfind"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/markbates/safe"
)

func byBuilder(f plugfind.Finder) plugfind.Finder {
	fn := func(name string, plugs []plugins.Plugin) plugins.Plugin {
		p := f.Find(name, plugs)
		if p == nil {
			return nil
		}
		if c, ok := p.(Builder); ok {
			if c.PluginName() == name {
				return p
			}
		}
		return nil
	}
	return plugfind.FinderFn(fn)
}

func (bc *Cmd) Main(ctx context.Context, root string, args []string) error {
	flags := bc.Flags()
	if err := flags.Parse(args); err != nil {
		return err
	}

	if len(flags.Args()) == 0 && bc.help {
		return plugprint.Print(plugio.Stdout(bc.ScopedPlugins()...), bc)
	}

	plugs := bc.ScopedPlugins()

	if len(flags.Args()) > 0 {
		name := flags.Args()[0]
		fn := plugfind.Background()
		fn = byBuilder(fn)
		fn = plugcmd.ByNamer(fn)
		fn = plugcmd.ByAliaser(fn)

		p := fn.Find(name, plugs)
		if p == nil {
			return fmt.Errorf("unknown builder %q", name)
		}
		b, ok := p.(Builder)
		if !ok {
			return fmt.Errorf("unknown builder %q", name)
		}
		return b.Build(ctx, root, args[1:])
	}

	info, err := here.Dir(root)
	if err != nil {
		return err
	}

	if info.Name != "main" {
		fp := filepath.Join(root, "cmd", info.Name)
		if _, err := os.Stat(fp); err == nil {
			info, err = here.Dir(fp)
			if err != nil {
				return err
			}
			root = fp
			os.Chdir(root)
		}
	}

	if err = bc.beforeBuild(ctx, root, args); err != nil {
		err = fmt.Errorf("before build %w", err)
		return bc.afterBuild(ctx, root, args, err)
	}

	if !bc.skipTemplateValidation {
		for _, p := range plugs {
			tv, ok := p.(TemplatesValidator)
			if !ok {
				continue
			}
			err = safe.RunE(func() error {
				return tv.ValidateTemplates(info.Dir)
			})
			if err != nil {
				return bc.afterBuild(ctx, root, args, err)
			}
		}
	}

	err = safe.RunE(func() error {
		err := bc.pack(ctx, info, plugs)
		if err != nil {
			err = fmt.Errorf("pack error %w", err)
		}
		return nil
	})
	if err != nil {
		return bc.afterBuild(ctx, root, args, err)
	}

	err = safe.RunE(func() error {
		err := bc.build(ctx, root)
		if err != nil {
			err = fmt.Errorf("build error %w", err)
		}
		return nil
	})

	return bc.afterBuild(ctx, root, args, err)

}

func (cmd *Cmd) beforeBuild(ctx context.Context, root string, args []string) error {
	plugs := cmd.ScopedPlugins()
	for _, p := range plugs {
		if bb, ok := p.(BeforeBuilder); ok {
			err := safe.RunE(func() error {
				err := bb.BeforeBuild(ctx, root, args)
				if err != nil {
					return fmt.Errorf("%s: %w", p.PluginName(), err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cmd *Cmd) afterBuild(ctx context.Context, root string, args []string, err error) error {
	plugs := cmd.ScopedPlugins()
	for _, p := range plugs {
		if bb, ok := p.(AfterBuilder); ok {
			err := safe.RunE(func() error {
				return bb.AfterBuild(ctx, root, args, err)
			})
			if err != nil {
				return err
			}
		}
	}
	return err
}
