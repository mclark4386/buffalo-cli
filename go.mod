module github.com/gobuffalo/buffalo-cli

go 1.13

replace github.com/gobuffalo/buffalo-cli/plugins => ./plugins

require (
	github.com/gobuffalo/attrs v0.1.0
	github.com/gobuffalo/buffalo-cli/plugins v0.0.0-00010101000000-000000000000
	github.com/gobuffalo/fizz v1.9.5
	github.com/gobuffalo/flect v0.2.0
	github.com/gobuffalo/genny/v2 v2.0.1
	github.com/gobuffalo/here v0.6.0
	github.com/gobuffalo/meta/v2 v2.0.0
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/pop/v5 v5.0.6
	github.com/markbates/grift v1.5.0
	github.com/markbates/jim v0.5.0
	github.com/markbates/pkger v0.14.0
	github.com/markbates/refresh v1.10.0
	github.com/markbates/safe v1.0.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/tools v0.0.0-20200117220505-0cba7a3a9ee9
)
