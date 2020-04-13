package build

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_MainFile_Version(t *testing.T) {
	r := require.New(t)

	bc := &MainFile{}

	ctx := context.Background()

	s, err := bc.Version(ctx, "")
	r.NoError(err)
	r.Contains(s, `"time":`)

	bc.WithPlugins(func() []plugins.Plugin {
		return plugins.Plugins{
			&buildBuilder{},
			&buildVersioner{version: "v1"},
		}
	})
	bc.HidePlugin()

	s, err = bc.Version(ctx, "")
	r.NoError(err)
	r.Contains(s, `"time":`)
	r.Contains(s, `"buildVersioner":"v1"`)
}

func Test_MainFile_Version_Err(t *testing.T) {
	r := require.New(t)

	bc := &MainFile{}

	ctx := context.Background()

	bc.WithPlugins(func() []plugins.Plugin {
		return plugins.Plugins{
			&buildBuilder{},
			&buildVersioner{version: ""},
			&buildVersioner{err: errors.New("Bad Version")},
		}
	})

	s, err := bc.Version(ctx, "")
	r.Error(err)
	r.Contains(err.Error(), "Bad Version")
	r.Equal("", s)
}

func Test_MainFile_generateNewMain(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	ctx := context.Background()
	bb := &bytes.Buffer{}
	err := bc.generateNewMain(ctx, ref, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.Contains(out, `appcli "github.com/markbates/coke/cli"`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `appcli.Buffalo`)
	r.Contains(out, `originalMain`)

	bc = &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}
	err = bc.generateNewMain(ctx, ref, "v1", bb)
	r.NoError(err)

	out = bb.String()
	r.Contains(out, `appcli "github.com/markbates/coke/cli"`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `appcli.Buffalo`)
	r.Contains(out, `originalMain`)
}

func Test_MainFile_generateNewMain_noCli(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return false },
	}

	bb := &bytes.Buffer{}
	err := bc.generateNewMain(context.Background(), ref, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.NotContains(out, `appcli "github.com/markbates/coke/cli"`)
	r.NotContains(out, `appcli.Buffalo`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `originalMain`)
	r.Contains(out, `cb.Main`)
}

func Test_MainFile_generateNewMain_Errs(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
			err: errors.New("Bad Import"),
		},
		&buildTagger{},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	ctx := context.Background()
	bb := &bytes.Buffer{}

	err := bc.generateNewMain(ctx, ref, "v1", bb)
	r.Error(err)
	r.Contains(err.Error(), "Bad Import")

	plugs = plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
		&buildTagger{},
	}
	bc = &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	//Bad writer
	bw := &badWriter{
		err: errors.New("Bad Writer"),
	}
	err = bc.generateNewMain(ctx, ref, "v1", bw)
	r.Error(err)
	r.Contains(err.Error(), "Bad Writer")
}

//TODO: Test BeforeBuild
func Test_MainFile_BeforeBuild(t *testing.T) {
	r := require.New(t)

	root := filepath.Join(".", "coke")
	println(root)
	current_root, err := os.Getwd()
	r.NoError(err)

	err, ref := setupProject(root, t)
	defer teardownProject(root, current_root)
	r.NoError(err)

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	ctx := context.Background()
	var args []string
	err = bc.BeforeBuild(ctx, root, args)
	r.NoError(err)
}

//TODO: Test AfterBuild
//TODO: Test renameMain

func setupProject(root string, t *testing.T) (error, here.Info) {
	if err := os.MkdirAll(root, os.ModeDir|os.ModePerm); err != nil {
		return err, here.Info{}
	}
	if err := os.Chdir(root); err != nil {
		return err, here.Info{}
	}

	//create files
	if err := ioutil.WriteFile("main.go", []byte(maingofile), 0660); err != nil {
		return err, here.Info{}
	}
	if err := ioutil.WriteFile("go.mod", []byte(gomodfile), 0660); err != nil {
		return err, here.Info{}
	}
	return nil, newRef(t, root)
}

func teardownProject(root, newRoot string) {
	os.Chdir(newRoot)
	os.RemoveAll(root)
}

const maingofile = `
package main

func main() {
	print("go time!")
}
`

const gomodfile = `
module github.com/markbates/coke

go 1.13`
