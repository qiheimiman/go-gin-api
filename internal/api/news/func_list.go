package news

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type listRequest struct{}

type listData struct {
	ID          int64     `json:"id" gorm:"id"`                     // 自增主键
	Title       string    `json:"title" gorm:"title"`               // 新闻标题
	SubTitle    string    `json:"sub_title" gorm:"sub_title"`       // 新闻副标题
	WebLink     string    `json:"web_link" gorm:"web_link"`         // 新闻网页链接
	BinanceId   string    `json:"binance_id" gorm:"binance_id"`     // Binance平台新闻ID
	PublishTime time.Time `json:"publish_time" gorm:"publish_time"` // 新闻发布时间
	Type        int8      `json:"type" gorm:"type"`                 // 类型 1-币安新闻
}

type listResponse struct {
	List []*listData `json:"list"`
}

// List 获取最新加密市场新闻列表
// @Summary 获取最新加密市场新闻列表
// @Description 获取最新加密市场新闻列表
// @Tags PUBLIC_API.news
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param Request body listRequest true "请求信息"
// @Success 200 {object} listResponse
// @Failure 400 {object} code.Failure
// @Router /public_api/news/list [get]
func (h *handler) List() core.HandlerFunc {
	return func(c core.Context) {
		info, err := h.newsService.List(c)
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
				Title:       v.Title,
				SubTitle:    v.SubTitle,
				PublishTime: v.PublishTime,
				WebLink:     v.WebLink,
				Type:        int8(v.Type),
			})
		}
		c.Payload(res)
	}
}
