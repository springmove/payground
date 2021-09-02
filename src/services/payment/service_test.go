package payment

import (
	"fmt"
	"testing"

	"github.com/linshenqi/payground/src/base"
)

func getSrv() *Service {
	srv := Service{
		cfg: Config{
			Endpoints: map[string]base.PaymentEndpoint{
				"alipay": {
					Provider: base.PaymentAlipay,
					MchKey:   "2088241146552860",
					// MchSecret: "binxb6prlj8c5nehsqy2mvaq95xsemx4",
					MchSecret: "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQGGTY/bO0fUyf/uwYtRZOLV+TMjCk3VbvUGbaTm/Eza3lVNtanB9+/rygg6oZ76psaG3tAcxSxY8BxXOhf3qBxVZYw2VWN0X5V24ggfLGDvuA/b29cyp0P6bFBJ64jQXzhVVy5F4YyO0vh3Ue7eMW4oPnqQhDUusHyGQ483eANwIDAQAB",
				},
			},
		},
	}

	if err := srv.initProviders(nil); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &srv
}

func TestService(t *testing.T) {
	srv := getSrv()
	if srv == nil {
		return
	}
}
