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
					Provider:  base.PaymentAlipay,
					MchKey:    "",
					MchSecret: "",
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

	resp, err := srv.CreatePayment("alipay", &base.Payment{
		Type:     base.PaymentTypeAlipayScan,
		TradeNo:  "awef123",
		TotalFee: 0.01,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}
