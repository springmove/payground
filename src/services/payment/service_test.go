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

	eid := "awef"
	resp, err := srv.CreatePayment("alipay", &base.Payment{
		Type:     base.PaymentTypeAlipayScan,
		TradeNo:  eid,
		TotalFee: 0.01,
		Desc:     "测试订单",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)

	respQuery, err := srv.GetPayment("alipay", &base.PaymentQuery{
		TradeNo: eid,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(respQuery)
}
