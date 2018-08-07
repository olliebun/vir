package main

import (
	"github.com/urfave/cli"

	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/ceralena/vir/virErrors"
)

func fatal(err error) {
	fmt.Println("error: " + err.Error())
	// XXX(cera) - is this appropriate for all platforms, including windows?
	syscall.Exit(1)
}

func main() {
	runCliApp(os.Args)
}

func runCliApp(args []string) {
	app := cli.NewApp()

	app.Name = "vir"
	app.Usage = "manage your music library"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "music-root, mr",
			Usage:  "music root directory",
			EnvVar: "VIR_MUSIC_ROOT",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "list-files",
			Aliases: []string{"ls"},
			Usage:   "list all music files",
			Action:  makeAction(actionListFiles),
		},
		{
			Name:    "rebuild-index",
			Aliases: []string{"r"},
			Usage:   "rebuild the vir index",
			Action:  makeAction(actionRebuildIndex),
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(args)

	if err != nil {
		fatal(err)
	}
}

type virAction func(*virContext, *cli.Context) virErrors.ScopedError

type virContext struct {
	musicLibraryRoot string
}

func makeAction(fn virAction) func(ctx *cli.Context) error {
	return func(cliCtx *cli.Context) error {
		virCtx := &virContext{
			cliCtx.GlobalString("music-root"),
		}
		err := fn(virCtx, cliCtx)
		return err
	}
}
