package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration of app
type Configuration struct {
	RedisURl string `required:"true" split_words:"true"`
	AcmeURl  string `required:"true" split_words:"true"`
	Domain   string `required:"true"`
	HTTPPort string `required:"true" split_words:"true"`
	UDPPort  string `required:"true" split_words:"true"`
	IPV4     string `required:"true"`
	IPV6     string
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
