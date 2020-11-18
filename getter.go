package loader

// Getter is implemented by types can self-serialize keys.
type Getter interface {
	// Merge return a new key with merge prefix and key
	Merge(prefix string, key string) string
	// Get return (value, found, error) with specified key
	Get(key string) (string, bool, error)
}
