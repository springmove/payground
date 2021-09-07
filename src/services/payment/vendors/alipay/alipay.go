package alipay

import (
	"fmt"

	"github.com/linshenqi/payground/src/base"
	v3 "github.com/smartwalle/alipay/v3"
)

func Payment2AlipayPreCreateReq(payment *base.Payment) *v3.TradePreCreate {

	req := v3.TradePreCreate{}
	req.OutTradeNo = payment.TradeNo
	req.TotalAmount = fmt.Sprintf("%.2f", payment.TotalFee)
	req.Subject = payment.Desc

	return &req
}

func Payment2AlipayTradeWapPayReq(payment *base.Payment) *v3.TradeWapPay {

	req := v3.TradeWapPay{}
	req.OutTradeNo = payment.TradeNo
	req.TotalAmount = fmt.Sprintf("%.2f", payment.TotalFee)
	req.Subject = payment.Desc

	return &req
}

func AlipayPreCreateResp2PaymentResp(resp *v3.TradePreCreateRsp) *base.CreatePaymentResp {

	paymentResp := base.CreatePaymentResp{
		PaymentUrl: resp.Content.QRCode,
	}

	return &paymentResp
}


func AlipayTradeWapPayResp2PaymentResp(resp *v3.TradePreCreateRsp) *base.CreatePaymentResp {

	paymentResp := base.CreatePaymentResp{
		PaymentUrl: resp.Content.QRCode,
	}

	return &paymentResp
}

func AlipayTradeQueryResp2PaymentNotify(resp *v3.TradeQueryRsp) *base.PaymentNotify {

	paymentNotify := base.PaymentNotify{
		TradeNo: resp.Content.OutTradeNo,
	}

	switch resp.Content.TradeStatus {
	case v3.TradeStatusWaitBuyerPay:
		paymentNotify.Status = base.PaymentStatusNotPay

	case v3.TradeStatusClosed:
		paymentNotify.Status = base.PaymentStatusClosed

	case v3.TradeStatusSuccess:
		paymentNotify.Status = base.PaymentStatusSuccess

	case v3.TradeStatusFinished:
		paymentNotify.Status = base.PaymentStatusFinished

	}
	return &paymentNotify
}
