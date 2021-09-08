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
					CertFile:  "",
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
		Type:     base.PaymentTypeAlipayWap,
		TradeNo:  "aw233ef234234",
		TotalFee: 0.01,
		Desc:     "测试订单",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)

	// resp, err := srv.GetPayment("alipay", &base.PaymentQuery{
	// 	TradeNo: "awef123",
	// })

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// fmt.Println(resp)
}
