package build

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build/buildtest"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_GoBuilder(t *testing.T) {
	cli := func(n ...string) string {
		cli := filepath.Join(n...)
		if runtime.GOOS == "windows" {
			cli += ".exe"
		}
		return cli
	}

	table := []struct {
		args []string
		exp  []string
	}{
		{exp: []string{"go", "build", "-o", cli("bin", "build")}},
		{
			args: []string{"-o", filepath.Join("bin", "foo")},
			exp:  []string{"go", "build", "-o", cli("bin", "foo")},
		},
		{
			args: []string{"-mod", "vendor"},
			exp:  []string{"go", "build", "-mod", "vendor", "-o", cli("bin", "build")},
		},
		{
			args: []string{"-tags", "a b c"},
			exp:  []string{"go", "build", "-tags", "abc", "-o", cli("bin", "build")},
		},
		{
			args: []string{"-static"},
			exp:  []string{"go", "build", "-o", cli("bin", "build"), "-ldflags", "-linkmode external -extldflags \"-static\""},
		},
		{
			args: []string{"-ldflags", "linky"},
			exp:  []string{"go", "build", "-o", cli("bin", "build"), "-ldflags", "-linkmode external -extldflags \"-static\""},
		},
	}

	for _, tt := range table {
		t.Run(strings.Join(tt.args, " "), func(st *testing.T) {
			r := require.New(st)

			bc := &Cmd{}

			var act []string
			fn := func(ctx context.Context, root string, args []string) error {
				act = args
				return nil
			}
			bc.WithPlugins(func() []plugins.Plugin {
				return []plugins.Plugin{
					buildtest.GoBuilder(fn),
				}
			})

			ctx := context.Background()
			err := bc.Main(ctx, "", tt.args)
			r.NoError(err)

			r.Equal(tt.exp, act)
		})
	}
}
