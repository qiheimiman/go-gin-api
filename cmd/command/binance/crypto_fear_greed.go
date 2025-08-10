package binance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/xinliangnote/go-gin-api/cmd/global"
	"github.com/xinliangnote/go-gin-api/cmd/initialize"
	"github.com/xinliangnote/go-gin-api/cmd/pkg/lock"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var BinanceFearGreedCmd = &cobra.Command{
	Use:   "binance_fear_greed",
	Short: "币安贪婪恐惧值",
	Long:  "币安贪婪恐惧值",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(context.Background(), "traceID", uuid.New().String())
		// 初始化配置
		initialize.InitConfig(global.Env)
		// 初始化日志
		initialize.SetupLogger("binance_fear_greed")
		initialize.InitErpDB(ctx)

		// 检查 global.DB 是否初始化成功
		if global.DB == nil {
			global.Logger.Error(ctx, "数据库未初始化，global.DB is nil")
			fmt.Println("错误：数据库未初始化")
			return // 或者 os.Exit(1)
		}

		lock, err := lock.GetFileLock("binance_fear_greed.lock")
		if err != nil {
			global.Logger.Error(ctx, "脚本正在执行,获取锁失败:"+err.Error())
			fmt.Printf("获取锁失败: %v\n", err)
			return
		}
		// 释放互斥锁
		defer lock.Release()
		global.Logger.Info(ctx, "开始执行脚本")

		// 目标 URL
		url := "https://www.binance.com/zh-CN/square/fear-and-greed-index"

		var htmlContent string

		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false), // 设置为非无头模式，显示浏览器界面
			// 对于Chrome 112+，可能需要添加这个
			chromedp.Flag("headless=new", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("start-maximized", true),
			chromedp.Flag("no-sandbox", true),
		)

		// 创建一个新的Chrome实例
		allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
		defer cancel()

		// 创建一个新的浏览器上下文
		browserCtx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		// 执行浏览器任务
		err = chromedp.Run(browserCtx,
			chromedp.Navigate(url),                   // 打开页面
			chromedp.Sleep(10*time.Second),           // 等待页面加载
			chromedp.OuterHTML("html", &htmlContent), // 获取页面 HTML
		)
		if err != nil {

			return
		}

		// 使用正则表达式提取 <script> 标签中 id="__APP_DATA" 的内容
		re := regexp.MustCompile(`<script id="__APP_DATA" type="application/json" nonce="[^"]*">(.+?)</script>`)
		matches := re.FindStringSubmatch(htmlContent)

		if len(matches) < 2 {
			global.Logger.Error(ctx, "未找到 __APP_DATA 的内容")
			return
		}

		// 提取到的 JSON 内容
		appData := matches[1]
		// fmt.Println("提取到的 __APP_DATA 内容:")
		// fmt.Println(appData)

		// 将提取到的内容写入文件
		filePath := "app_data.json"
		err = os.WriteFile(filePath, []byte(appData), 0644)
		if err != nil {
			fmt.Println("写入文件失败:", err)
			return
		}

		// 定义一个通用的结构来解析 JSON
		var data map[string]interface{}

		// 解析 JSON 数据
		err = json.Unmarshal([]byte(appData), &data)
		if err != nil {
			return
		}

		// 提取 "dataByRouteId" 字段
		loader := data["appState"].(map[string]interface{})["loader"].(map[string]interface{})
		dataByRouteId := loader["dataByRouteId"].(map[string]interface{})

		// 遍历 "dataByRouteId" 找到包含 "fearGreed" 的键
		var fearGreed map[string]interface{}
		for _, value := range dataByRouteId {
			entry := value.(map[string]interface{})
			if fg, ok := entry["fearGreed"]; ok {
				fearGreed = fg.(map[string]interface{})
				break
			}
		}

		// 检查是否找到 "fearGreed"
		if fearGreed == nil {
			err = errors.New("未找到 fearGreed 字段")
			return
		}

		// 打印 "fearGreed" 的内容
		global.Logger.Infof(ctx, "fearGreed 字段内容:")
		global.Logger.Infof(ctx, "当前值 (currentValue): %.0f\n", fearGreed["currentValue"])
		global.Logger.Infof(ctx, "昨日值 (yesterdayValue): %.0f\n", fearGreed["yesterdayValue"])
		global.Logger.Infof(ctx, "上周值 (lastWeekValue): %.0f\n", fearGreed["lastWeekValue"])
		global.Logger.Infof(ctx, "看跌票数 (bearishValue): %.0f\n", fearGreed["bearishValue"])
		global.Logger.Infof(ctx, "看涨票数 (bullishValue): %.0f\n", fearGreed["bullishValue"])

		// 获取今天的时间
		today := time.Now()
		// 获取昨天的时间
		// yesterday := today.AddDate(0, 0, -1)
		// 定义日期格式
		layout := "2006-01-02"

		// 格式化日期
		todayStr := today.Format(layout)
		// yesterdayStr := yesterday.Format(layout)

		var oldData struct {
			ID int64
		}

		// 更新今天贪婪恐惧值
		result := global.DB.Table("crypto_fear_greed").Where("date", todayStr).Limit(1).Find(&oldData)
		if result.Error != nil {
			return
		}

		if result.RowsAffected > 0 { // 存在数据
			result = global.DB.Table("crypto_fear_greed").Where("date", todayStr).Updates(map[string]interface{}{
				"date":            todayStr,
				"current_value":   int8(fearGreed["currentValue"].(float64)),
				"bearish_value":   int(fearGreed["bearishValue"].(float64)),
				"bullish_value":   int(fearGreed["bullishValue"].(float64)),
				"last_week_value": int8(fearGreed["lastWeekValue"].(float64)),
			})

		} else {
			result = global.DB.Table("crypto_fear_greed").Create(map[string]interface{}{
				"date":            todayStr,
				"current_value":   int8(fearGreed["currentValue"].(float64)),
				"bearish_value":   int(fearGreed["bearishValue"].(float64)),
				"bullish_value":   int(fearGreed["bullishValue"].(float64)),
				"last_week_value": int8(fearGreed["lastWeekValue"].(float64)),
			})
		}

		if result.Error != nil {
			global.Logger.Error(ctx, "更新greed_fear数据失败:"+result.Error.Error())
		}

	},
}

func init() {
	// 环境参数
	BinanceFearGreedCmd.Flags().StringVarP(&global.Env, "env", "e", "dev", "测试环境")
}
