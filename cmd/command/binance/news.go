package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/xinliangnote/go-gin-api/cmd/global"
	"github.com/xinliangnote/go-gin-api/cmd/initialize"
	"github.com/xinliangnote/go-gin-api/cmd/pkg/lock"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/news"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var BinanceNewsCmd = &cobra.Command{
	Use:   "binance_news",
	Short: "币安新闻",
	Long:  "币安新闻",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(context.Background(), "traceID", uuid.New().String())
		// 初始化配置
		initialize.InitConfig(global.Env)
		// 初始化日志
		initialize.SetupLogger("binance_news")
		initialize.InitErpDB(ctx)

		// 检查 global.DB 是否初始化成功
		if global.DB == nil {
			global.Logger.Error(ctx, "数据库未初始化，global.DB is nil")
			fmt.Println("错误：数据库未初始化")
			return // 或者 os.Exit(1)
		}

		lock, err := lock.GetFileLock("binance_news.lock")
		if err != nil {
			global.Logger.Error(ctx, "脚本正在执行,获取锁失败:"+err.Error())
			fmt.Printf("获取锁失败: %v\n", err)
			return
		}
		// 释放互斥锁
		defer lock.Release()

		// 设置超时时间
		ctx, _ = context.WithTimeout(ctx, 15*time.Second)

		// 目标 URL
		url := "https://www.binance.com/zh-CN/square/news/all"
		var htmlContent string

		// 执行浏览器任务
		// 创建Chrome浏览器选项
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			// chromedp.Flag("headless", false), // 设置为非无头模式，显示浏览器界面
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

		err = chromedp.Run(browserCtx,
			chromedp.Evaluate(`navigator.webdriver = undefined`, nil), // 隐藏 webdriver
			chromedp.Navigate(url),                   // 打开页面
			chromedp.Sleep(3*time.Second),            // 等待页面加载
			chromedp.OuterHTML("html", &htmlContent), // 获取页面 HTML
		)
		if err != nil {
			global.Logger.Error(ctx, err)
		}

		// 打印抓取的 HTML 内容
		// fmt.Println("页面内容已抓取..." + htmlContent)

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
		// // fmt.Println(appData)

		// // 将提取到的内容写入文件
		// filePath := "app_data.json"
		// err = os.WriteFile(filePath, []byte(appData), 0644)
		// if err != nil {
		// 	global.Logger.Error(ctx, "写入文件失败", err)
		// }

		// 定义一个通用的结构来解析 JSON
		var data map[string]interface{}

		// 解析 JSON 数据
		err = json.Unmarshal([]byte(appData), &data)
		if err != nil {
			global.Logger.Error(ctx, "JSON 解析失败", err)
			return
		}

		// 提取 "dataByRouteId" 字段
		loader := data["appState"].(map[string]interface{})["loader"].(map[string]interface{})
		dataByRouteId := loader["dataByRouteId"].(map[string]interface{})

		// 遍历 "dataByRouteId" 找到包含 "fearGreed" 的键
		for _, value := range dataByRouteId {
			entry := value.(map[string]interface{})

			if fg, ok := entry["list"]; ok {
				// fmt.Println(fg)
				list := fg.([]interface{})

				for _, value2 := range list {
					row := value2.(map[string]interface{})

					var news news.News

					// 是否存在该新闻
					result := global.DB.Table("news").Where("binance_id = ?", row["id"].(string)).Limit(1).Find(&news)
					if result.Error != nil {
						global.Logger.Error(ctx, result.Error)
						return
					}
					// 已存在该新闻
					if result.RowsAffected > 0 {
						continue
					}

					// 新增
					result = global.DB.Table("news").Create(map[string]interface{}{
						"title":     row["title"].(string),
						"sub_title": row["subTitle"].(string),
						"web_link":  row["webLink"].(string),

						"binance_id":   row["id"].(string),
						"publish_time": time.Unix(int64(row["date"].(float64)), 0).Format("2006-01-02 15:04:05"),
					})

					if result.Error != nil {
						global.Logger.Error(ctx, result.Error)
					}
				}

			}
		}
	},
}

func init() {
	// 环境参数
	BinanceNewsCmd.Flags().StringVarP(&global.Env, "env", "e", "dev", "测试环境")
}
