package mugglepay

import "github.com/linshenqi/payground/src/base"

type Config struct {
	Token string `yaml:"token"`
}

func (s *Config) ConfigName() string {
	return base.ServiceMugglePay
}

func (s *Config) Validate() error {
	return nil
}

func (s *Config) Default() interface{} {
	return &Config{}
}
