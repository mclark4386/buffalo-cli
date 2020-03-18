package clifix

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cmd_Fix(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)

	f, err := os.Create(filepath.Join(dir, "go.mod"))
	r.NoError(err)
	f.WriteString("module coke")
	r.NoError(f.Close())

	ctx := context.Background()
	var args []string

	fixer := &Cmd{}
	err = fixer.Fix(ctx, dir, args)
	r.NoError(err)

	root := filepath.Join(dir, "cmd", "buffalo")

	_, err = os.Stat(root)
	r.NoError(err)

	fp := filepath.Join(root, "main.go")
	b, err := ioutil.ReadFile(fp)
	r.NoError(err)
	r.Contains(string(b), `"coke/cmd/buffalo/cli"`)

	_, err = os.Stat(filepath.Join(root, "cli", "buffalo.go"))
	r.NoError(err)
}
