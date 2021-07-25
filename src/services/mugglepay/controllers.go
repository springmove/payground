package mugglepay

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/payground/src/base"
	"github.com/linshenqi/sptty"
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

	fmt.Println(req)

END:
	_ = sptty.SimpleResponse(ctx, status, map[string]int{
		"status": status,
	})
}
