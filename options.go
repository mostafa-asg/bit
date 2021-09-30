package bit

// Options need for set initialization
type Options struct {
	// number of initial bits
	nbits int
}

type Option func(*Options)

// WithInitialBits sets the initial bits required for the set
func WithInitialBits(n int) Option {
	return func(opts *Options) {
		opts.nbits = n
	}
}
