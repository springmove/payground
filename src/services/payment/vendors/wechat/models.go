package wechat

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/springmove/payground/src/base"
)

const (
	NotifyEndpoint = "/v1/payment-wechat"
)

const (
	TradeTypeMiniProgram = "JSAPI"
)

const (
	ResultSuccess = "SUCCESS"
	ResultFail    = "FAIL"
)

func generateSign(payload interface{}, secret string) string {
	vals := map[string]interface{}{}
	body, _ := json.Marshal(payload)
	_ = json.Unmarshal(body, &vals)

	keys := []string{}
	for k := range vals {
		if k == "sign" || k == "XMLName" {
			continue
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)

	str := ""
	for _, v := range keys {
		str += fmt.Sprintf("%s=%v&", v, vals[v])
	}

	str += fmt.Sprintf("key=%s", secret)

	h := md5.New()
	_, _ = h.Write([]byte(str))

	sign := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	return sign
}

type PayloadBase struct {
	XMLName xml.Name `xml:"xml"`
}

type ReqOrder struct {
	PayloadBase

	AppKey         string `xml:"appid" json:"appid"`
	MchKey         string `xml:"mch_id" json:"mch_id"`
	NonceStr       string `xml:"nonce_str" json:"nonce_str"`
	Sign           string `xml:"sign" json:"sign"`
	Boby           string `xml:"body" json:"body"`
	TradeNo        string `xml:"out_trade_no" json:"out_trade_no"`
	TotalFee       int    `xml:"total_fee" json:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip" json:"spbill_create_ip"`
	NotifyUrl      string `xml:"notify_url" json:"notify_url"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	OpenID         string `xml:"openid" json:"openid"`
}

func (s *ReqOrder) FromPayment(payment *base.Payment) {
	s.AppKey = payment.AppKey
	s.TotalFee = int(payment.TotalFee)
	s.TradeNo = payment.TradeNo
	s.Boby = payment.Desc

	switch payment.Type {
	case base.PaymentTypeWechatMiniProgram:
		s.TradeType = TradeTypeMiniProgram

	default:
		s.TradeType = TradeTypeMiniProgram
	}
}

func (s *ReqOrder) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespErr struct {
	ErrCode     string `xml:"err_code" json:"err_code"`
	ErrCodeDesc string `xml:"err_code_des" json:"err_code_des"`
}

type RespReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

type ResqOrder struct {
	PayloadBase
	RespReturn
	RespErr
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
}

type ReqNotify struct {
	PayloadBase
	RespReturn
	ResultCode string `xml:"result_code"`
	TradeNo    string `xml:"out_trade_no"`
	TradeState string `xml:"trade_state"`
}

func (s *ReqNotify) ToPaymentNotify() *base.PaymentNotify {
	notify := base.PaymentNotify{
		TradeNo: s.TradeNo,
		Msg:     fmt.Sprintf("%+v", s.RespReturn),
	}

	switch s.TradeState {
	case "SUCCESS":
		notify.Status = base.PaymentStatusSuccess

	case "REFUND":
		notify.Status = base.PaymentStatusRefund

	case "NOTPAY":
		notify.Status = base.PaymentStatusNotPay

	case "CLOSED":
		notify.Status = base.PaymentStatusClosed

	case "REVOKED":
		notify.Status = base.PaymentStatusRevoked

	case "USERPAYING":
		notify.Status = base.PaymentStatusPaying

	default:
		notify.Status = base.PaymentStatusUnKnown
	}

	return &notify
}

type RespNotify struct {
	PayloadBase
	RespReturn
}

type ReqOrderQuery struct {
	PayloadBase

	AppKey   string `xml:"appid" json:"appid"`
	MchKey   string `xml:"mch_id" json:"mch_id"`
	TradeNo  string `xml:"out_trade_no" json:"out_trade_no"`
	NonceStr string `xml:"nonce_str" json:"nonce_str"`
	Sign     string `xml:"sign" json:"sign"`
}

func (s *ReqOrderQuery) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type ReqTransferQuery struct {
	PayloadBase

	AppKey   string `xml:"appid" json:"appid"`
	MchKey   string `xml:"mch_id" json:"mch_id"`
	NonceStr string `xml:"nonce_str" json:"nonce_str"`
	Sign     string `xml:"sign" json:"sign"`
	TradeNo  string `xml:"partner_trade_no" json:"partner_trade_no"`
}

func (s *ReqTransferQuery) FromTransferQuery(transfer *base.QueryTransfer) {
	s.AppKey = transfer.AppKey
	s.TradeNo = transfer.TradeNo
}

func (s *ReqTransferQuery) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespQueryTransfer struct {
	PayloadBase

	RespReturn
	RespErr

	ResultCode string `xml:"result_code"`
	Status     string `xml:"status"`
	Reason     string `xml:"reason"`
}

func (s *RespQueryTransfer) ToQueryTransferResp() *base.QueryTransferResp {
	resp := base.QueryTransferResp{
		Reason: s.Reason,
	}

	switch s.Status {
	case "SUCCESS":
		resp.Status = base.TransferStatusSuccess
	case "FAILED":
		resp.Status = base.TransferStatusFailed
	case "PROCESSING":
		resp.Status = base.TransferStatusProcessing

	default:
		resp.Status = base.TransferStatusProcessing
	}

	return &resp
}

type ReqTransfer struct {
	PayloadBase

	AppKey   string `xml:"mch_appid" json:"mch_appid"`
	MchKey   string `xml:"mchid" json:"mchid"`
	NonceStr string `xml:"nonce_str" json:"nonce_str"`
	Sign     string `xml:"sign" json:"sign"`
	TradeNo  string `xml:"partner_trade_no" json:"partner_trade_no"`

	OpenID    string `xml:"openid" json:"openid"`
	CheckName string `xml:"check_name" json:"check_name"`
	TotalFee  int    `xml:"amount" json:"amount"`
	Desc      string `xml:"desc" json:"desc"`
}

func (s *ReqTransfer) FromPaymentTransfer(transfer *base.PaymentTransfer) {
	s.AppKey = transfer.AppKey
	s.Desc = transfer.Desc
	s.TotalFee = transfer.TotalFee
	s.TradeNo = transfer.TradeNo
	s.OpenID = transfer.OpenID
}

func (s *ReqTransfer) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespTransfer struct {
	PayloadBase

	RespReturn
	RespErr

	ResultCode string `xml:"result_code"`
	TradeNo    string `xml:"partner_trade_no"`
	PaymentNo  string `xml:"payment_no"`
}

type ReqRefund struct {
	PayloadBase

	AppKey    string `xml:"appid" json:"appid"`
	MchKey    string `xml:"mch_id" json:"mch_id"`
	NonceStr  string `xml:"nonce_str" json:"nonce_str"`
	Sign      string `xml:"sign" json:"sign"`
	TradeNo   string `xml:"out_trade_no" json:"out_trade_no"`
	RefundNo  string `xml:"out_refund_no" json:"out_refund_no"`
	TotalFee  int    `xml:"total_fee" json:"total_fee"`
	RefundFee int    `xml:"refund_fee" json:"refund_fee"`
}

func (s *ReqRefund) FromPayment(payment *base.Payment) {
	s.AppKey = payment.AppKey
	s.TradeNo = payment.TradeNo
	s.RefundNo = payment.TradeNo
	s.TotalFee = int(payment.TotalFee)

	if payment.RefundFee > 0 {
		s.RefundFee = payment.RefundFee
	} else {
		s.RefundFee = int(payment.TotalFee)
	}
}

func (s *ReqRefund) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespRefund struct {
	PayloadBase

	RespReturn
	RespErr

	ResultCode string `xml:"result_code"`
	TradeNo    string `xml:"out_trade_no" json:"out_trade_no"`
}

type ReqQueryRefund struct {
	PayloadBase

	AppKey   string `xml:"appid" json:"appid"`
	MchKey   string `xml:"mch_id" json:"mch_id"`
	NonceStr string `xml:"nonce_str" json:"nonce_str"`
	Sign     string `xml:"sign" json:"sign"`
	TradeNo  string `xml:"out_trade_no" json:"out_trade_no"`
}

func (s *ReqQueryRefund) FromQueryRefund(query *base.QueryRefund) {
	s.AppKey = query.AppKey
	s.TradeNo = query.TradeNo
}

func (s *ReqQueryRefund) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespQueryRefund struct {
	PayloadBase

	RespReturn
	RespErr

	ResultCode string `xml:"result_code"`
	TradeNo    string `xml:"out_trade_no" json:"out_trade_no"`

	Status string `xml:"refund_status_0"`
}

func (s *RespQueryRefund) ToQueryRefundResp() *base.QueryRefundResp {
	resp := base.QueryRefundResp{}

	switch s.Status {
	case "SUCCESS":
		resp.Status = base.RefundStatusSuccess
	case "REFUNDCLOSE":
		resp.Status = base.RefundStatusClosed
	case "PROCESSING":
		resp.Status = base.RefundStatusProcessing
	case "CHANGE":
		resp.Status = base.RefundStatusException

	default:
		resp.Status = base.TransferStatusProcessing
	}

	return &resp
}

type ReqClosePayment struct {
	PayloadBase

	AppKey   string `xml:"appid" json:"appid"`
	MchKey   string `xml:"mch_id" json:"mch_id"`
	NonceStr string `xml:"nonce_str" json:"nonce_str"`
	Sign     string `xml:"sign" json:"sign"`
	TradeNo  string `xml:"out_trade_no" json:"out_trade_no"`
}

func (s *ReqClosePayment) FromPayment(payment *base.Payment) {
	s.AppKey = payment.AppKey
	s.TradeNo = payment.TradeNo
}

func (s *ReqClosePayment) GenerateSign(secret string) {
	s.Sign = generateSign(s, secret)
}

type RespClosePayment struct {
	PayloadBase

	RespReturn
	RespErr

	ResultCode string `xml:"result_code"`
}
