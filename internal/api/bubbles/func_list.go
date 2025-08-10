package bubbles

import (
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type listRequest struct{}

type listData struct {
	ID          int64   `json:"id"`          // 自增主键
	Rank        int     `json:"rank"`        // 货币排名
	Price       float64 `json:"price"`       // price
	Volume      int64   `json:"volume"`      // 24小时交易量
	Marketcap   int64   `json:"marketcap"`   // 市值
	Performance string  `json:"performance"` // 新闻发布时间

}

type listResponse struct {
	List []*listData `json:"list"`
}

// List 获取最新排名前1000加密货币价格
// @Summary 获取最新排名前1000加密货币价格
// @Description 获取最新排名前1000加密货币价格
// @Tags PUBLIC_API.bubbles
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param Request body listRequest true "请求信息"
// @Success 200 {object} listResponse
// @Failure 400 {object} code.Failure
// @Router /public_api/bubbles/list [get]
func (h *handler) List() core.HandlerFunc {
	return func(c core.Context) {
		info, err := h.bubblesService.List(c)
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

		// 初始化一个新的 listResponse 实例
		res := &listResponse{
			List: make([]*listData, 0, len(info)),
		}
		// 将info赋值给res
		// 将 info 中的元素复制到 res 指向的切片中
		for _, v := range info {
			res.List = append(res.List, &listData{
				ID:          int64(v.Id),
				Rank:        int(v.Rank),
				Price:       v.Price,
				Volume:      v.Volume,
				Marketcap:   v.Marketcap,
				Performance: v.Performance,
			})
		}
		c.Payload(res)
	}
}
