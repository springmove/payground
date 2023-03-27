package mugglepay

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/springmove/payground/src/base"
	"github.com/springmove/sptty"
)

func (s *Service) routePostMugglePayCallBack(ctx iris.Context) {
	ctx.Header("Content-Type", "application/json")

	req := base.MuggleRespOrder{}
	var err error
	status := iris.StatusOK

	if err = ctx.ReadJSON(&req); err != nil {
		sptty.Log(sptty.ErrorLevel, fmt.Sprintf("routePostMugglePayCallBack Failed: %s", err.Error()), s.ServiceName())
		status = iris.StatusBadRequest
		goto END
	}

	if err = s.serviceDispatcher.Dispatch(base.DispatcherMuggleCallback, req); err != nil {
		sptty.Log(sptty.ErrorLevel, fmt.Sprintf("routePostMugglePayCallBack.Dispatch Failed: %s", err.Error()), s.ServiceName())
	}

END:
	_ = sptty.SimpleResponse(ctx, status, map[string]int{
		"status": status,
	})
}
