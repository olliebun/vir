// Package index provides the core of vir
// it provides indexing capabilities and
// data structures for a music library
package index

import (
	"github.com/ceralena/vir/state"
	"github.com/ceralena/vir/track"
	"github.com/ceralena/vir/util"
	"github.com/ceralena/vir/virErrors"

	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Index represents a Vir music library index.
type Index interface {
	Rebuild() virErrors.ScopedError

	// Yield a full list of music files.
	ListMusicFiles() <-chan MusicFileListEntry
}

// MusicFileListEntry is a single entry from ListMusicFiles.
// It can contain either a RelFilename or an Error.
type MusicFileListEntry struct {
	RelFilename string
	Error       virErrors.ScopedError
}

// LoadIndex loads a vir index from a given root directory.
func LoadIndex(musicLibraryRoot string) (Index, virErrors.ScopedError) {
	// check that the music dir actually exists
	err := checkMusicLibraryRootExists(musicLibraryRoot)
	if err != nil {
		return nil, err
	}

	// initialize a state cache
	stateCache, err := state.GetStateCache()
	if err != nil {
		return nil, err
	}

	// return an Index interface
	return &index{musicLibraryRootDir: musicLibraryRoot, stateCache: stateCache}, nil
}

type index struct {
	musicLibraryRootDir string
	stateCache          state.Cache
}

func (idx *index) getFullPath(relPath string) string {
	return filepath.Join(idx.musicLibraryRootDir, relPath)
}

func (idx *index) Rebuild() virErrors.ScopedError {
	err := idx.stateCache.Set("musicRoot", []byte(idx.musicLibraryRootDir))

	if err != nil {
		return err
	}

	files := idx.ListMusicFiles()

	for fileEntry := range files {
		if fileEntry.Error != nil {
			// FIXME(cera)- we shouldn't just give up here
			// we'll actually get an orphaned goroutine if we do!
			return fileEntry.Error
		}

		tr, err := track.LoadTrackFromPath(idx.getFullPath(fileEntry.RelFilename))
		if err != nil {
			// FIXME(cera)- we shouldn't just give up here
			// we'll actually get an orphaned goroutine if we do!
			return err
		}

		// XXX(cera) - do something else here
		fmt.Printf("%#v\n", tr)
	}

	return nil
}

func stripRootDirFromPath(rootDir, path string) string {
	strippedPath := strings.Replace(path, rootDir, "", 1)
	if strippedPath[0] == filepath.Separator {
		return strippedPath[1:]
	}
	return strippedPath
}

func (idx *index) ListMusicFiles() <-chan MusicFileListEntry {
	ch := make(chan MusicFileListEntry)

	go func() {
		defer close(ch)
		walkFunc := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			// strip the path of the root prefix
			// XXX(cera) - do we actually want to strip the root dir here, if we're just putting it back
			// whenever we have to read or write with the file?
			strippedPath := stripRootDirFromPath(idx.musicLibraryRootDir, path)

			// send it to the channel
			ch <- MusicFileListEntry{strippedPath, nil}
			return nil
		}

		err := filepath.Walk(idx.musicLibraryRootDir, walkFunc)
		if err != nil {
			ch <- MusicFileListEntry{
				RelFilename: "",
				Error:       virErrors.ErrMusicLibraryWalkError("vir/index.ListMusicFiles", err),
			}
		}
	}()

	return ch
}

func checkMusicLibraryRootExists(musicLibraryRoot string) virErrors.ScopedError {
	exists, err := util.DirExists(musicLibraryRoot)
	if err != nil && util.IsPathIsNotDir(err) {
		return virErrors.ErrMusicLibraryRootIsNotDir("vir/index.checkMusicLibraryRootExists", err)
	} else if err != nil {
		return virErrors.ErrFatal("vir/index.checkMusicLibraryRootExists", err)
	} else if !exists {
		return virErrors.ErrMusicLibraryRootDoesNotExist("vir/index.checkMusicLibraryRootExists", musicLibraryRoot)
	}
	return nil
}
