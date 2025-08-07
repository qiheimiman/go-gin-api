package helper

import (
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type nowTimeResponse struct {
	NowTime string `json:"now_time"` // 当前时间
}

// NowTime 当前时间
// @Summary 当前时间
// @Description 当前时间
// @Tags Helper
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} nowTimeResponse
// @Failure 400 {object} code.Failure
// @Router /helper/now_time [get]
func (h *handler) NowTime() core.HandlerFunc {
	return func(ctx core.Context) {

		res := new(nowTimeResponse)
		res.NowTime = time.Now().Format("2006-01-02 15:04:05")
		ctx.Payload(res)
	}
}
