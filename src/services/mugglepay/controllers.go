package mugglepay

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/sptty"
)

func (s *Service) routePostMugglePayCallBack(ctx iris.Context) {
	ctx.Header("Content-Type", "application/json")

	// todo

	_ = sptty.SimpleResponse(ctx, http.StatusOK, map[string]int{
		"status": http.StatusOK,
	})
}
