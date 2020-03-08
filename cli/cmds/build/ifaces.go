package build

import (
	"context"
	"flag"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/spf13/pflag"
)

// Builder is a sub-command of buffalo build.
// 	buffalo build webpack
type Builder interface {
	plugins.Plugin
	Build(ctx context.Context, root string, args []string) error
}

type BeforeBuilder interface {
	plugins.Plugin
	BeforeBuild(ctx context.Context, root string, args []string) error
}

type AfterBuilder interface {
	plugins.Plugin
	AfterBuild(ctx context.Context, root string, args []string, err error) error
}

type Flagger interface {
	plugins.Plugin
	BuildFlags() []*flag.Flag
}

type Pflagger interface {
	plugins.Plugin
	BuildFlags() []*pflag.Flag
}

type TemplatesValidator interface {
	plugins.Plugin
	ValidateTemplates(root string) error
}

type Packager interface {
	plugins.Plugin
	Package(ctx context.Context, root string, files []string) error
}

type PackFiler interface {
	plugins.Plugin
	PackageFiles(ctx context.Context, root string) ([]string, error)
}

type Versioner interface {
	plugins.Plugin
	BuildVersion(ctx context.Context, root string) (string, error)
}

type Importer interface {
	plugins.Plugin
	BuildImports(ctx context.Context, root string) ([]string, error)
}

type GoBuilder interface {
	// GoBuild will be called to build, and execute, the
	// presented context and args.
	// The first plugin to receive this call will be the
	// only to answer it.
	GoBuild(ctx context.Context, root string, args []string) error
}

type BuildArger interface {
	// GoBuildArgs receives the current list
	// and returns either the same list, or
	// a modified version of the arguments.
	// Implementations are responsible for ensuring
	// that the arguments returned are "valid"
	// arguments for the `go build` command.
	GoBuildArgs(args []string) []string
}

type Stdouter = plugio.Outer
type Stdiner = plugio.Inner
type Stderrer = plugio.Errer
