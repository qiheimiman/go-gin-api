package zhonghuan

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/cmd/global"
	"github.com/xinliangnote/go-gin-api/cmd/initialize"
	"github.com/xinliangnote/go-gin-api/cmd/pkg/lock"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var username string
var password string

var ZhonghuanCmd = &cobra.Command{
	Use:   "zhonghuan",
	Short: "中环转运签到",
	Long:  "中环转运签到",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(context.Background(), "traceID", uuid.New().String())
		// 初始化配置
		initialize.InitConfig(global.Env)
		// 初始化日志
		initialize.SetupLogger("zhonghuan")
		// initialize.InitErpDB(ctx)
		lock, err := lock.GetFileLock("zhonghuan.lock")
		if err != nil {
			global.Logger.Error(ctx, "脚本正在执行,获取锁失败:"+err.Error())
			return
		}
		// 释放互斥锁
		defer lock.Release()

		// 创建一个HTTP客户端，带有Cookie支持
		jar, _ := cookiejar.New(nil)
		client := &http.Client{Jar: jar}

		loginURL := "https://xinzhongus.com/webLogin"
		data2 := `{"username":"` + username + `","password":"` + password + `","loginType":1}`

		req, err := http.NewRequest("POST", loginURL, strings.NewReader(data2))

		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// 读取响应体
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		type apiResponse struct {
			Code       int         `json:"code"`
			Msg        string      `json:"msg"`
			Data       interface{} `json:"data"`
			Count      int         `json:"count"`
			CurrentDay int         `json:"currentDay"`
			Status     string      `json:"status"`
		}
		var apiResp apiResponse
		if err = json.Unmarshal(responseBody, &apiResp); err != nil {
			return
		}

		// 构建后续JSON POST请求
		postURL := "https://xinzhongus.com/integral/onclickSign"

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		dateParam := time.Now().Format("2006-01-02")
		signDate := strings.ReplaceAll(dateParam, "-", "")
		// 添加普通表单字段
		writer.WriteField("date", signDate)
		writer.Close()

		// 创建请求对象
		req1, err := http.NewRequest("POST", postURL, &requestBody)
		if err != nil {
			return
		}

		req1.Header.Set("Content-Type", writer.FormDataContentType())
		req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
		// 发送后续POST请求
		resp1, err := client.Do(req1)
		if err != nil {
			return
		}
		defer resp1.Body.Close()

		// 读取并处理响应
		responseBody, err = io.ReadAll(resp1.Body)
		if err != nil {
			global.Logger.Error(ctx, "Error reading response body", zap.Error(err))
			return
		}
		global.Logger.Info(ctx, "Response Body:", zap.String("body", string(responseBody)))

	},
}

func init() {
	// 环境参数
	ZhonghuanCmd.Flags().StringVarP(&global.Env, "env", "e", "dev", "测试环境")
	// 用户名参数
	ZhonghuanCmd.Flags().StringVarP(&username, "username", "u", "", "登录用户名")
	// 密码参数
	ZhonghuanCmd.Flags().StringVarP(&password, "password", "p", "", "登录密码")
	// 是否开启调试模式
}
