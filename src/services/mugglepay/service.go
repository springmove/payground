package mugglepay

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/springmove/payground/src/base"
	"github.com/springmove/sptty"
	"gopkg.in/resty.v1"

	medusaBase "github.com/springmove/medusa/src/base"
)

type Service struct {
	sptty.BaseService

	http *resty.Client
	cfg  Config

	serviceDispatcher medusaBase.IServiceDispatcher
}

func (s *Service) Init(app sptty.ISptty) error {

	s.http = sptty.CreateHttpClient(sptty.DefaultHttpClientConfig())

	if app == nil {
		return nil
	}

	if err := app.GetConfig(s.ServiceName(), &s.cfg); err != nil {
		return err
	}

	s.serviceDispatcher = app.GetService(medusaBase.ServiceDispatcher).(medusaBase.IServiceDispatcher)
	if s.serviceDispatcher == nil {
		return fmt.Errorf("%s Service Is Required", medusaBase.ServiceDispatcher)
	}

	app.AddRoute("POST", "/v1/mugglepay-callback", s.routePostMugglePayCallBack)

	if err := s.serviceDispatcher.CreateDispatcher(base.DispatcherMuggleCallback); err != nil {
		return err
	}

	return nil
}

func (s *Service) ServiceName() string {
	return base.ServiceMugglePay
}

func (s *Service) SetToken(token string) {
	s.cfg.Token = token
}

// 创建订单
func (s *Service) CreateOrder(req *base.MuggleReqOrder) (*base.MuggleRespOrderCreate, error) {
	url := "https://api.mugglepay.com/v1/orders"

	r := s.http.R().
		SetBody(req).
		SetHeader("content-type", "application/json").
		SetHeader(base.MuggleHeaderToken, s.cfg.Token)

	resp, err := r.Post(url)
	if err != nil {
		return nil, err
	}

	respOrderCreate := base.MuggleRespOrderCreate{}
	if err := json.Unmarshal(resp.Body(), &respOrderCreate); err != nil {
		return nil, err
	}

	if respOrderCreate.Status != http.StatusCreated {
		return nil, fmt.Errorf("%+v", respOrderCreate.MuggleRespBase)
	}

	return &respOrderCreate, nil
}

// 查询订单
func (s *Service) GetOrder(orderID string) (*base.MuggleRespOrderCreate, error) {
	url := fmt.Sprintf("https://api.mugglepay.com/v1/orders/%s", orderID)

	r := s.http.R().
		SetHeader("content-type", "application/json").
		SetHeader(base.MuggleHeaderToken, s.cfg.Token)

	resp, err := r.Get(url)
	if err != nil {
		return nil, err
	}

	respOrderCreate := base.MuggleRespOrderCreate{}
	if err := json.Unmarshal(resp.Body(), &respOrderCreate); err != nil {
		return nil, err
	}

	if respOrderCreate.Status != http.StatusOK {
		return nil, fmt.Errorf("%+v", respOrderCreate.MuggleRespBase)
	}

	return &respOrderCreate, nil
}

// 取消订单
func (s *Service) CancelOrder() {

}

// 退款
func (s *Service) Refund() {

}
