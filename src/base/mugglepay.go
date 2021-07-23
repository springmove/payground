package base

const (
	ServiceMugglePay = "mugglepay"
)

const (
	MuggleHeaderToken = "token"
)

type MuggleReqOrder struct {
	// 用户订单id
	MerchantOrderID string `json:"merchant_order_id"`

	// 必填
	PriceAmount float64 `json:"price_amount"`

	// 必填
	PriceCurrency string `json:"price_currency"`

	// 支付渠道
	PayCurrency string `json:"pay_currency"`

	// 支付结果回调
	CallbackUrl string `json:"callback_url"`

	// 取消支付跳转链接
	CancelUrl string `json:"cancel_url"`

	// 支付成功跳转链接
	SuccessUrl string `json:"success_url"`

	// 用户支付结果回调时的验证token
	Token string `json:"token"`

	// 支付标题
	Title string `json:"title"`

	// 支付描述
	Description string `json:"description"`

	// Based on PC or Mobile Wap, we provide different links, for Alipay / Alipay Global / Wechat only.
	Mobile bool `json:"mobile"`

	// Return the payment url directly, for Alipay / Alipay Global / Wechat only.
	Fast bool `json:"fast"`
}

type MuggleRespOrder struct {
	MuggleReqOrder
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type MuggleRespBase struct {
	Status    int    `json:"status"`
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type MuggleRespOrderCreate struct {
	MuggleRespBase

	PaymentUrl string           `json:"payment_url"`
	Order      *MuggleRespOrder `json:"order"`
}
