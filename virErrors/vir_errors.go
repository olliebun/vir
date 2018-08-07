// Package virErrors provides all errors produced by vir code.
package virErrors

import (
	"fmt"
)

func scopedErr(scope, msg string) ScopedError {
	return scopedError{scope, msg}
}

// ScopedError is an error type with a Scope() method.
// Its key purpose is to ensure that all calls up and down the chain within vir have been wrapped with scope context.
type ScopedError interface {
	Error() string
	Scope() string
}

type scopedError struct {
	scope string
	msg   string
}

func (sErr scopedError) Error() string {
	return fmt.Sprintf("%s: %s", sErr.scope, sErr.msg)
}

func (sErr scopedError) Scope() string {
	return sErr.scope
}

// ErrNotImplemented is used to stub out code where we haven't implemented some feature.
func ErrNotImplemented(scope, feature string) ScopedError {
	return scopedErr(scope, "NotImplemented: "+feature)
}

// ErrUserLookupFailed is used when we can't look up a user to get the home dir for storing vir state.
func ErrUserLookupFailed(scope string, err error) ScopedError {
	return scopedErr(scope, "could not look up your user; this means vir can't configure itself: "+err.Error())
}

// ErrCacheSetupFailed is used when we fail to set up a cache for configuration.
func ErrCacheSetupFailed(scope string, err error) ScopedError {
	return scopedErr(scope, "could not set up configuration cache: "+err.Error())
}

// ErrCacheOperationFailed is used when a cache operation fails.
func ErrCacheOperationFailed(scope, op, key string, err error) ScopedError {
	return scopedErr(scope, fmt.Sprintf("cache operation %s for key %s failed: %s", op, key, err))
}

// ErrMusicLibraryRootDoesNotExist is used when vir is asked to index a music library root directory that does not exist.
func ErrMusicLibraryRootDoesNotExist(scope, path string) ScopedError {
	return scopedErr(scope, "music library root dir does not exist: "+path)
}

// ErrMusicLibraryRootIsNotDir is used when the music library root path specified to vir exists but is not a directory.
func ErrMusicLibraryRootIsNotDir(scope string, err error) ScopedError {
	return scopedErr(scope, err.Error())
}

// ErrMusicLibraryWalkError is used when vir encounters an error while walking the music library.
func ErrMusicLibraryWalkError(scope string, err error) ScopedError {
	return scopedErr(scope, "encountered an error while walking the music library: "+err.Error())
}

// ErrTrackID3MetadataLoadFailed is used when we fail to read id3 metadata for a file.
func ErrTrackID3MetadataLoadFailed(scope, fullPath string, err error) ScopedError {
	return scopedErr(scope, fmt.Sprintf("encountered an error while loading id3 metadata for %s: %s", fullPath, err))
}

// ErrFatal is used when we encounter an unexpected I/O error or some other kind of fatal error that is very difficult
// to predict or recover from.
//
// Use sparingly. Overuse is a code smell and will lead to opaque error states.
func ErrFatal(scope string, err error) ScopedError {
	return scopedErr(scope, fmt.Sprintf("unexpected fatal error: (%T): %s", err, err))
}
