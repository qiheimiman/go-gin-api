package helper

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type html2urlRequest struct {
	HtmlCode string `form:"html_code" binding:"required"` // 需要转换的HTML代码
}

type html2urlResponse struct {
	Url string `json:"url"` // 网页链接
}

// html2url 将html代码转为可访问的url链接
// @Summary 将html代码转为可访问的url链接
// @Description 将html代码转为可访问的url链接
// @Tags Helper
// @Accept application/json
// @Produce json
// @Success 200 {object} nowTimeResponse
// @Failure 400 {object} code.Failure
// @Router /helper/html2url [post]
func (h *handler) Html2url() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(html2urlRequest)
		res := new(html2urlResponse)

		if err := ctx.ShouldBindForm(req); err != nil {
			fmt.Println(err)
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		fileName := "html2url_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".html"
		// 生成随机文件名
		filePath := fmt.Sprintf("/www/wwwroot/resource.abcd1234.top/%v", fileName)
		// filePath := fmt.Sprintf("/Users/xjh/Downloads/%v", fileName)

		// 移除HTML代码中的头尾标签
		htmlCode := strings.TrimPrefix(strings.TrimSuffix(req.HtmlCode, "```"), "```html")
		// 将HTML内容写入文件
		err := os.WriteFile(filePath, []byte(htmlCode), 0644)
		if err != nil {
			fmt.Println(err)
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(err),
			)
			return
		}

		res.Url = "http://resource.abcd1234.top/" + fileName
		ctx.Payload(res)
	}
}
