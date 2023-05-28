package env

import (
	"os"
	"strings"
)

// Getter is implemented by types can self-serialize keys.
type Getter interface {
	// Merge return a new key with merge prefix and key
	Merge(prefix string, key string) string

	// Get return (value, found, error) with specified key
	Get(key string) (string, bool, error)
}

type getter struct{}

func (g *getter) Merge(prefix string, key string) string {
	var nk string // new key
	if prefix != "" && key != "" {
		nk = prefix + "_" + key
	} else if prefix != "" {
		nk = prefix
	} else {
		nk = key
	}
	return nk
}

func (g *getter) Get(key string) (string, bool, error) {
	key = strings.ToUpper(key)
	value, found := os.LookupEnv(key)
	return value, found, nil
}
