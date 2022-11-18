package kong

import "github.com/alecthomas/kong"

// Run a cli using kong
func Run(cli interface{}) {
	ctx := kong.Parse(cli, kong.UsageOnError())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
