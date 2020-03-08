package build

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/markbates/safe"
)

func (bc *Cmd) buildArgs(root string) ([]string, error) {
	args := []string{"build"}

	info, err := here.Dir(root)
	if err != nil {
		return nil, err
	}

	bin := bc.bin
	if len(bin) == 0 {
		bin = filepath.Join("bin", info.Name)
	}

	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") {
			bin += ".exe"
		}
		bin = strings.Replace(bin, "/", "\\", -1)
	} else {
		bin = strings.TrimSuffix(bin, ".exe")
	}
	args = append(args, "-o", bin)

	if len(bc.mod) != 0 {
		args = append(args, "-mod", bc.mod)
	}

	args = append(args, bc.buildFlags...)

	flags := []string{}

	if bc.static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(bc.ldFlags) > 0 {
		flags = append(flags, bc.ldFlags)
	}

	if len(flags) > 0 {
		args = append(args, "-ldflags", strings.Join(flags, " "))
	}

	for _, p := range bc.ScopedPlugins() {
		if bt, ok := p.(BuildArger); ok {
			args = bt.GoBuildArgs(args)
		}
	}

	return args, nil
}

func (bc *Cmd) build(ctx context.Context, root string) error {
	buildArgs, err := bc.buildArgs(root)
	if err != nil {
		return err
	}

	plugs := bc.ScopedPlugins()
	for _, p := range plugs {
		if br, ok := p.(GoBuilder); ok {
			return safe.RunE(func() error {
				return br.GoBuild(ctx, root, buildArgs)
			})
		}
	}

	cmd := exec.CommandContext(ctx, "go", buildArgs...)
	cmd.Stdin = plugio.Stdin(plugs...)
	cmd.Stdout = plugio.Stdout(plugs...)
	cmd.Stderr = plugio.Stderr(plugs...)
	return cmd.Run()
}
