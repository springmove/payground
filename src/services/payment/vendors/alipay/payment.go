package alipay

import (
	"fmt"

	"github.com/linshenqi/payground/src/base"
	v3 "github.com/smartwalle/alipay/v3"
)

type PaymentProvider struct {
	base.BasePaymentProvider

	client *v3.Client
}

func (s *PaymentProvider) Init(paymentUrl string, endpoint *base.PaymentEndpoint) error {
	client, err := v3.New(endpoint.MchKey, endpoint.MchSecret, true)
	if err != nil {
		return fmt.Errorf("Create Alipay Client Failed: %s", err.Error())
	}

	// if err := client.LoadAliPayPublicCert(endpoint.CertFile); err != nil {
	// 	return fmt.Errorf("Load PublicKey Failed: %s", err.Error())
	// }

	s.client = client

	return s.BasePaymentProvider.Init(paymentUrl, endpoint)
}

func (s *PaymentProvider) CreatePayment(payment *base.Payment) (*base.CreatePaymentResp, error) {

	var paymentResp *base.CreatePaymentResp

	switch payment.Type {
	case base.PaymentTypeAlipayScan:
		resp, err := s.client.TradePreCreate(*Payment2AlipayPreCreateReq(payment))
		if err != nil {
			return nil, err
		}

		if resp.Content.Code != v3.CodeSuccess {
			return nil, fmt.Errorf("Alipay.CreatePayment Failed Code:%s Msg:%s", resp.Content.Code, resp.Content.Msg)
		}

		paymentResp = AlipayPreCreateResp2PaymentResp(resp)

	default:
		return nil, fmt.Errorf("Payment Type Not Supported: %s", payment.Type)
	}

	paymentResp.Type = payment.Type
	return paymentResp, nil
}

func (s *PaymentProvider) GetPayment(query *base.PaymentQuery) (*base.PaymentNotify, error) {
	var notify *base.PaymentNotify

	resp, err := s.client.TradeQuery(v3.TradeQuery{
		OutTradeNo: query.TradeNo,
	})

	if err != nil {
		return nil, err
	}

	if resp.Content.Code != v3.CodeSuccess {
		return nil, fmt.Errorf("Alipay.GetPayment Failed Code:%s Msg:%s", resp.Content.Code, resp.Content.Msg)
	}

	notify = AlipayTradeQueryResp2PaymentNotify(resp)

	return notify, nil
}

func (s *PaymentProvider) Refund(payment *base.Payment) error {
	return nil
}

func (s *PaymentProvider) QueryRefund(query *base.QueryRefund) (*base.QueryRefundResp, error) {

	return nil, nil
}
