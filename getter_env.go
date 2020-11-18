package loader

import (
	"os"
	"strings"
)

type EnvGetter struct{}

func (g *EnvGetter) Merge(prefix string, key string) string {
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

func (g *EnvGetter) Get(key string) (string, bool, error) {
	key = strings.ToUpper(key)
	value, found := os.LookupEnv(key)
	return value, found, nil
}
