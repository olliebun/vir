package track

import (
	"fmt"

	"github.com/ceralena/vir/virErrors"

	"github.com/casept/id3-go"
	"strconv"
	"strings"
)

// Metadata represents the metadata of a track from its id3 tag.
type Metadata struct {
	Title  string
	Artist string
	Album  string
	Number int
}

// Track represents a track with the metadata read and parsed from its id3 tag.
type Track struct {
	FullPath string
	Metadata
	Errata []string
}

// LoadTrackFromPath loads a track from its id3 data given the full path to a file.
func LoadTrackFromPath(fullPath string) (*Track, virErrors.ScopedError) {
	var errata []string

	f, err := id3.Open(fullPath)

	if err != nil {
		return nil, virErrors.ErrTrackID3MetadataLoadFailed("vir/track.LoadTrackFromPath", fullPath, err)
	}

	defer func() {
		e := f.Close()
		if e != nil {
			// FIXME(cera) - what can we do here?
			panic(err)
		}
	}()

	trackNum, err := parseTrackNumber(f)

	if err != nil {
		// we consider this a warning, not an error
		errata = append(errata, err.Error())
	}

	metadata := Metadata{
		Title:  clean(f.Tagger.Title()),
		Artist: clean(f.Tagger.Artist()),
		Album:  clean(f.Tagger.Album()),
		Number: trackNum,
	}

	return &Track{
		FullPath: fullPath,
		Metadata: metadata,
		Errata:   errata,
	}, nil

}

func parseTrackNumber(tagger id3.Tagger) (int, error) {
	trackFrame := tagger.Frame("TRCK")
	if trackFrame == nil {
		return -1, fmt.Errorf("no track number (TRCK) in tags")
	}

	// get a clean frame string
	frameStr := clean(trackFrame.String())

	// now check for forward slashes
	// this is to handle a syntax like this:
	// 1/8
	if strings.ContainsRune(frameStr, '/') {
		spl := strings.SplitN(frameStr, "/", 2)
		if len(spl) != 2 {
			return -1, fmt.Errorf("invalid or unrecognised track number format in tags: %s", frameStr)
		}
		frameStr = spl[0]
	}

	trackNum, err := strconv.Atoi(frameStr)

	if err != nil {
		return -1, fmt.Errorf("invalid track number format in tags: %s: %s", frameStr, err)
	}

	return trackNum, nil

}

func clean(elem string) string {
	a := strings.Replace(elem, "\u0000", "", -1)
	b := strings.Replace(a, "\u0026", "", -1)
	c := strings.Replace(b, "\x00", "", -1)
	return c
}
