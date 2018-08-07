// Package util provides some simple helpers for vir code.
package util

import (
	"os"
)

type errPathIsNotDir struct {
	dirPath string
}

func (err errPathIsNotDir) Error() string {
	return "path exists but is not a dir: " + err.dirPath
}

// IsPathIsNotDir returns a boolean indicating whether the error is known to report that a path passed to DirExists was
// found to exist, but is not a directory (e.g. it is a regular file or a special file).
func IsPathIsNotDir(err error) bool {
	_, ok := err.(errPathIsNotDir)
	return ok
}

// DirExists reports on whether a directory exists.
//
// If there is an IO error when checking, an error is returned.
// If the path exists and it is not a directory, an error is returned.
//
// Otherwise, error is nil and the boolean value is whether the directory exists.
func DirExists(dirPath string) (bool, error) {
	fi, err := os.Stat(dirPath)

	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if !fi.IsDir() {
		return false, errPathIsNotDir{dirPath: dirPath}
	}
	return true, nil
}
