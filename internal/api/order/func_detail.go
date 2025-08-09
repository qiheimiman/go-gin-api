package order

import (
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/order"
)

type detailRequest struct {
	Id string `uri:"id"` // HashID
}

type detailResponse struct {
	Id      int32  `json:"id"`
	OrderNo string `json:"order_no"`
}

// Detail 查看详情
// @Summary 查看详情
// @Description 查看详情
// @Tags API.order
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param Request body detailRequest true "请求信息"
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Router /order/{id} [get]
func (h *handler) Detail() core.HandlerFunc {
	return func(c core.Context) {
		req := new(detailRequest)
		res := new(detailResponse)
		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// ids, err := h.hashids.HashidsDecode(req.Id)
		// if err != nil {
		// 	c.AbortWithError(core.Error(
		// 		http.StatusBadRequest,
		// 		code.HashIdsDecodeError,
		// 		code.Text(code.HashIdsDecodeError)).WithError(err),
		// 	)
		// 	return
		// }

		id := int32(1)

		searchOneData := new(order.SearchOneData)
		searchOneData.Id = id

		info, err := h.orderService.Detail(c, searchOneData)
		if err != nil {
			// 这里的逻辑是正确的，应该返回错误
			c.AbortWithError(core.Error(
				http.StatusInternalServerError, // 服务器错误更合适
				code.OrderDetailError,
				code.Text(code.OrderDetailError)).WithError(err),
			)
			return // 必须确保这里返回
		}
		// 只有在 err == nil 时才执行到这里
		if info == nil {
			// 额外防御性检查，虽然服务层应避免返回 nil, nil
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.OrderDetailError,
				code.Text(code.OrderDetailError)).WithError(fmt.Errorf("service returned nil info without error")),
			)
			return
		}

		res.Id = info.Id
		res.OrderNo = info.OrderNo
		fmt.Println(res)
		c.Payload(res)
	}
}
