package app

// Core application
type Core struct {
	Config   Config   `config:"-"`
	Resolver Resolver `resolver:"-"`
}

// New application core initialization
func New(c Config) *Core {
	return &Core{
		Config: c,
	}
}
