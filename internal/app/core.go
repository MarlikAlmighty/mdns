package app

// Core application
type Core struct {
	Config   Config   `config:"-"`
	Resolver Resolver `resolver:"-"`
	Store    Store    `store:"-"`
}

// New application core initialization
func New(c Config, r Resolver, s Store) *Core {
	return &Core{
		Config:   c,
		Resolver: r,
		Store:    s,
	}
}
