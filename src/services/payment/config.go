package payment

import (
	"github.com/linshenqi/payground/src/base"
	"github.com/linshenqi/sptty"
)

type Config struct {
	sptty.BaseConfig

	PaymentUrl string                          `yaml:"payment_url"`
	Endpoints  map[string]base.PaymentEndpoint `yaml:"endpoints"`
}

func (s *Config) ConfigName() string {
	return base.ServicePayment
}
