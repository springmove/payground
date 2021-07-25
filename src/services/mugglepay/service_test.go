package mugglepay

import (
	"fmt"
	"testing"

	"github.com/linshenqi/payground/src/base"
)

func getSrv() *Service {
	srv := Service{
		cfg: Config{
			Token: "",
		},
	}

	_ = srv.Init(nil)

	return &srv
}
func TestCreateOrder(t *testing.T) {

	srv := getSrv()
	resp, err := srv.CreateOrder(&base.MuggleReqOrder{
		MerchantOrderID: "test8",
		PriceAmount:     1,
		PriceCurrency:   "USD",
		PayCurrency:     base.MugglePayCurrencyAlipay,
		Token:           "1234",
		CallbackUrl:     "http://c2d730d2e636.ngrok.io/payground/v1/mugglepay-callback",
		// Fast:            true,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp.PaymentUrl)
	fmt.Println(*resp.Order)
}
