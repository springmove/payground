package wechat

import (
	"encoding/xml"
	"fmt"

	"github.com/springmove/payground/src/base"
	"github.com/springmove/sptty"
	"gopkg.in/resty.v1"
)

type PaymentProvider struct {
	base.BasePaymentProvider
	http *resty.Client
}

func (s *PaymentProvider) Init(paymentUrl string, endpoint *base.PaymentEndpoint) error {
	s.http = sptty.CreateHttpClient(sptty.DefaultHttpClientConfig())
	return s.BasePaymentProvider.Init(paymentUrl, endpoint)
}

func (s *PaymentProvider) CreatePayment(payment *base.Payment) (*base.CreatePaymentResp, error) {
	reqOrder := ReqOrder{
		MchKey:         s.BasePaymentProvider.Endpoint.MchKey,
		NonceStr:       sptty.GenerateUID(),
		SpbillCreateIP: s.BasePaymentProvider.GetReqIP(),
		NotifyUrl:      s.getNotifyUrl(),
		OpenID:         payment.OpenID,
	}

	reqOrder.FromPayment(payment)
	reqOrder.GenerateSign(s.Endpoint.MchSecret)

	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	body, _ := xml.Marshal(reqOrder)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	respBody := ResqOrder{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		return nil, fmt.Errorf("%+v", respBody)
	}

	return s.generatePaymentResp(payment.Type, respBody.PrepayID, reqOrder.NonceStr, payment.AppKey), nil
}

func (s *PaymentProvider) GetPayment(query *base.PaymentQuery) (*base.PaymentNotify, error) {
	req := ReqOrderQuery{
		AppKey:   query.AppKey,
		MchKey:   s.BasePaymentProvider.Endpoint.MchKey,
		NonceStr: sptty.GenerateUID(),
		TradeNo:  query.TradeNo,
	}

	req.GenerateSign(s.Endpoint.MchSecret)

	url := "https://api.mch.weixin.qq.com/pay/orderquery"

	body, _ := xml.Marshal(req)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	respBody := ReqNotify{}
	_ = xml.Unmarshal(resp.Body(), &respBody)

	return respBody.ToPaymentNotify(), nil
}

func (s *PaymentProvider) Transfer(transfer *base.PaymentTransfer) error {
	cert, err := s.loadCert()
	if err != nil {
		return err
	}

	s.http.SetCertificates(*cert)

	req := ReqTransfer{
		CheckName: "NO_CHECK",
	}

	req.MchKey = s.Endpoint.MchKey
	req.NonceStr = sptty.GenerateUID()

	req.FromPaymentTransfer(transfer)
	req.GenerateSign(s.Endpoint.MchSecret)

	url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	body, _ := xml.Marshal(req)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return err
	}

	respBody := RespTransfer{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		return fmt.Errorf("%+v", respBody)
	}

	return nil
}

func (s *PaymentProvider) QueryTransfer(query *base.QueryTransfer) (*base.QueryTransferResp, error) {
	cert, err := s.loadCert()
	if err != nil {
		return nil, err
	}

	s.http.SetCertificates(*cert)

	req := ReqTransferQuery{
		MchKey:   s.Endpoint.MchKey,
		NonceStr: sptty.GenerateUID(),
	}

	req.FromTransferQuery(query)
	req.GenerateSign(s.Endpoint.MchSecret)

	url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"
	body, _ := xml.Marshal(req)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	respBody := RespQueryTransfer{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		return nil, fmt.Errorf("%+v", respBody)
	}

	return respBody.ToQueryTransferResp(), nil
}

// 退款
func (s *PaymentProvider) Refund(payment *base.Payment) error {
	cert, err := s.loadCert()
	if err != nil {
		return err
	}

	s.http.SetCertificates(*cert)

	req := ReqRefund{
		MchKey:   s.Endpoint.MchKey,
		NonceStr: sptty.GenerateUID(),
	}

	req.FromPayment(payment)
	req.GenerateSign(s.Endpoint.MchSecret)

	url := "https://api.mch.weixin.qq.com/secapi/pay/refund"
	body, _ := xml.Marshal(req)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return err
	}

	respBody := RespRefund{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		return fmt.Errorf("%+v", respBody)
	}

	return nil
}

// 退款查询
func (s *PaymentProvider) QueryRefund(query *base.QueryRefund) (*base.QueryRefundResp, error) {
	req := ReqQueryRefund{
		MchKey:   s.Endpoint.MchKey,
		NonceStr: sptty.GenerateUID(),
	}

	req.FromQueryRefund(query)
	req.GenerateSign(s.Endpoint.MchSecret)

	url := "https://api.mch.weixin.qq.com/pay/refundquery"
	body, _ := xml.Marshal(req)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return nil, err
	}

	respBody := RespQueryRefund{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		return nil, fmt.Errorf("%+v", respBody)
	}

	return respBody.ToQueryRefundResp(), nil
}

// 关闭订单
func (s *PaymentProvider) Close(payment *base.Payment) error {
	req := ReqClosePayment{
		MchKey:   s.Endpoint.MchKey,
		NonceStr: sptty.GenerateUID(),
	}

	req.FromPayment(payment)
	req.GenerateSign(s.Endpoint.MchSecret)

	url := "https://api.mch.weixin.qq.com/pay/closeorder"
	body, _ := xml.Marshal(req)
	resp, err := s.http.R().SetBody(body).Post(url)
	if err != nil {
		return err
	}

	respBody := RespRefund{}
	_ = xml.Unmarshal(resp.Body(), &respBody)
	if respBody.ResultCode != ResultSuccess || respBody.ReturnCode != ResultSuccess {
		return fmt.Errorf("%+v", respBody)
	}

	return nil
}

func (s *PaymentProvider) GetNotifyController() *base.PaymentNotifyController {
	return &base.PaymentNotifyController{
		Method:   "POST",
		Endpoint: NotifyEndpoint,
		Handler:  s.notifyController,
	}
}
