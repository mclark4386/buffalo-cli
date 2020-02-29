package garlic

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/markbates/safe"
)

func Run(ctx context.Context, root string, args []string) error {
	main := filepath.Join(root, "cmd", "buffalo")
	if _, err := os.Stat(filepath.Dir(main)); err != nil {
		buff, err := cli.NewFromRoot(root)
		if err != nil {
			return err
		}
		return buff.Main(ctx, root, args)
	}

	bargs := []string{"run", "./cmd/buffalo"}
	bargs = append(bargs, args...)

	cmd := exec.CommandContext(ctx, "go", bargs...)
	cmd.Stdin = plugio.Stdin()
	cmd.Stdout = plugio.Stdout()
	cmd.Stderr = plugio.Stderr()
	err := safe.RunE(func() error {
		return cmd.Run()
	})
	if err != nil {
		return err
	}

	return nil
}
