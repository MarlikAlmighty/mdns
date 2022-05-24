package app

import (
	"github.com/MarlikAlmighty/mdns/internal/config"
)

// Core application
type Core struct {
	Resolver Resolver `resolver:"-"`
	Config   Config   `config:"-"`
}

// New application core initialization
func New(r Resolver, c *config.Configuration) *Core {
	return &Core{
		Resolver: r,
		Config:   c,
	}
}
