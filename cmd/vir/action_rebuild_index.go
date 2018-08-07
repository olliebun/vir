package main

import (
	"github.com/urfave/cli"

	"github.com/ceralena/vir/index"
	"github.com/ceralena/vir/virErrors"
)

// actionRebuildIndex is the CLI action for rebuild-index
func actionRebuildIndex(ctx *virContext, _ *cli.Context) virErrors.ScopedError {
	idx, err := index.LoadIndex(ctx.musicLibraryRoot)

	if err != nil {
		return err
	}

	err = idx.Rebuild()
	return err
}
