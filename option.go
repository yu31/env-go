package env

type options struct {
	prefix   string
	tagName  string
	override bool
	getter   Getter
}

type Option func(opts *options)

// WithTagName set struct's tag name
func WithTagName(name string) Option {
	return func(opts *options) {
		opts.tagName = name
	}
}

// WithPrefix set key's prefix
func WithPrefix(prefix string) Option {
	return func(opts *options) {
		opts.prefix = prefix
	}
}

// WithGetter set Getter instances
func WithGetter(getter Getter) Option {
	return func(opts *options) {
		opts.getter = getter
	}
}

// WithOverride in force override an existing value
func WithOverride(ok bool) Option {
	return func(opts *options) {
		opts.override = ok
	}
}
