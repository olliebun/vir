package main

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/ceralena/vir/index"
	"github.com/ceralena/vir/virErrors"
)

func actionListFiles(ctx *virContext, _ *cli.Context) virErrors.ScopedError {
	idx, err := index.LoadIndex(ctx.musicLibraryRoot)

	if err != nil {
		return err
	}

	fileEntries := idx.ListMusicFiles()

	for fileEntry := range fileEntries {
		if fileEntry.Error != nil {
			fatal(fileEntry.Error)
			break
		}
		fmt.Println(fileEntry.RelFilename)
	}

	return nil
}
