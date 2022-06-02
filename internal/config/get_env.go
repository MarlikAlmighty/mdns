package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration of app
type Configuration struct {
	HTTPPort    string   `required:"true" split_words:"true"`
	DnsHost     string   `required:"true" split_words:"true"`
	NameServers []string `required:"true" split_words:"true"`
	IPV6        bool     `required:"true"`
}

func New() *Configuration {
	return &Configuration{}
}

// GetEnv configuration init
func (cnf *Configuration) GetEnv() error {
	if err := envconfig.Process("", cnf); err != nil {
		return err
	}
	return nil
}
