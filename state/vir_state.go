// Package state provides facilities for vir to store state
package state

import (
	"github.com/ceralena/go-cacheh"

	"os/user"
	"path/filepath"

	"github.com/ceralena/vir/virErrors"
	"strings"
)

const stateVersion = "v1.0.0"

const virConfRootDir = ".vir"

func getCacheKeyPrefix() string {
	return strings.Replace(stateVersion, ".", "-", -1) + "-"
}

// getVirHome provides the root configuration directory for vir state.
func getVirConfRoot() (string, virErrors.ScopedError) {
	// look up the current user
	currentUser, err := user.Current()
	if err != nil {
		return "", virErrors.ErrUserLookupFailed("vir/state", err)
	}

	// find the user's home dir
	home := currentUser.HomeDir

	// join it to .vir
	return filepath.Join(home, virConfRootDir), nil
}

// Cache provides a simple persistent caching interface.
type Cache interface {
	Get(key string) ([]byte, virErrors.ScopedError)
	Set(key string, value []byte) virErrors.ScopedError
	Delete(key string) virErrors.ScopedError
}

type cache struct {
	cacheh.Cache
}

func (c *cache) Get(key string) ([]byte, virErrors.ScopedError) {
	val, err := c.Cache.Get(key)

	if err != nil {
		return nil, virErrors.ErrCacheOperationFailed("vir/state.Cache", "Get", key, err)
	}

	return val, nil

}

func (c *cache) Set(key string, value []byte) virErrors.ScopedError {
	err := c.Cache.Set(key, value)

	if err != nil {
		return virErrors.ErrCacheOperationFailed("vir/state.Cache", "Set", key, err)
	}

	return nil
}

func (c *cache) Delete(key string) virErrors.ScopedError {
	err := c.Cache.Delete(key)

	if err != nil {
		return virErrors.ErrCacheOperationFailed("vir/state.Cache", "Delete", key, err)
	}

	return nil
}

// GetStateCache provides a consistent state cache for vir state.
func GetStateCache() (Cache, virErrors.ScopedError) {
	confRoot, err := getVirConfRoot()
	if err != nil {
		return nil, err
	}

	cacheDsn := "dir:" + confRoot

	c, cacheErr := cacheh.NewCache(cacheDsn)

	if cacheErr != nil {
		return nil, virErrors.ErrCacheSetupFailed("vir/state:GetStateCache()", cacheErr)
	}

	// scope our cache keys to be prefixed by the music library root dir
	// FIXME(cera) - do not hard-code forward slash as filesep here
	return &cache{c.WithKeyPrefix(getCacheKeyPrefix())}, nil
}
