package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration of app
type Configuration struct {
	HTTPHost    string   `required:"true" split_words:"true"`
	HTTPPort    string   `required:"true" split_words:"true"`
	DnsHost     string   `required:"true" split_words:"true"`
	DnsTcpPort  string   `required:"true" split_words:"true"`
	DnsUdpPort  string   `required:"true" split_words:"true"`
	NameServers []string `required:"true" split_words:"true"`
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
