package bubbles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/xinliangnote/go-gin-api/cmd/global"
	"github.com/xinliangnote/go-gin-api/cmd/initialize"
	"github.com/xinliangnote/go-gin-api/cmd/pkg/lock"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/bubbles1000"

	"gorm.io/gorm"
)

// bubbles1000Data 结构体用于解析从 API 获取的 JSON 数据
// 字段名与 JSON key 匹配，并使用正确的 json tag
type bubbles1000Data struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Slug           string             `json:"slug"`
	Symbol         string             `json:"symbol"`
	Dominance      float64            `json:"dominance"`
	Image          string             `json:"image"`
	Rank           int64              `json:"rank"`
	Stable         bool               `json:"stable"` // 添加了缺失的字段
	Price          float64            `json:"price"`
	Marketcap      int64              `json:"marketcap"`
	Volume         int64              `json:"volume"`
	CgId           string             `json:"cg_id"` // 注意下划线
	Symbols        map[string]string  `json:"symbols"`
	Performance    map[string]float64 `json:"performance"`
	RankDiffs      map[string]int64   `json:"rank_diffs"`      // 注意下划线
	ExchangePrices map[string]float64 `json:"exchange_prices"` // 注意下划线
}

var Bubbles1000Cmd = &cobra.Command{
	Use:   "bubbles1000",
	Short: "获取加密泡泡前1000币价",
	Long:  "获取加密泡泡前1000币价",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(context.Background(), "traceID", uuid.New().String())
		// 初始化配置
		initialize.InitConfig(global.Env)
		// 初始化日志
		initialize.SetupLogger("bubbles1000")
		// *** 关键：初始化数据库 ***
		// 请根据你的项目实际情况调用正确的初始化函数，例如 initialize.InitDB(ctx)
		// 你原来的 initialize.InitErpDB(ctx) 可能不正确或不存在
		// initialize.InitErpDB(ctx)
		// 假设正确的初始化函数是 InitDB
		initialize.InitErpDB(ctx) // <--- 请确认这是你项目中正确的数据库初始化函数

		// 检查 global.DB 是否初始化成功
		if global.DB == nil {
			global.Logger.Error(ctx, "数据库未初始化，global.DB is nil")
			fmt.Println("错误：数据库未初始化")
			return // 或者 os.Exit(1)
		}

		lock, err := lock.GetFileLock("bubbles1000.lock")
		if err != nil {
			global.Logger.Error(ctx, "脚本正在执行,获取锁失败:"+err.Error())
			fmt.Printf("获取锁失败: %v\n", err)
			return
		}
		// 释放互斥锁
		defer lock.Release()

		// 修正 URL，移除末尾空格
		url := "https://cryptobubbles.net/backend/data/bubbles1000.usd.json" // *** 修正了 URL ***

		// 创建带超时的 context
		httpCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		// 发送GET请求
		req, err := http.NewRequestWithContext(httpCtx, "GET", url, nil)
		if err != nil {
			global.Logger.Error(ctx, "创建HTTP请求失败:"+err.Error())
			fmt.Printf("创建HTTP请求失败: %v\n", err)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			global.Logger.Error(ctx, "发送HTTP请求失败:"+err.Error())
			fmt.Printf("发送HTTP请求失败: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// 检查 HTTP 响应状态码
		if resp.StatusCode != http.StatusOK {
			global.Logger.Error(ctx, fmt.Sprintf("HTTP请求失败，状态码: %d", resp.StatusCode))
			fmt.Printf("HTTP 请求失败，状态码：%d\n", resp.StatusCode)
			return
		}

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			global.Logger.Error(ctx, "读取响应体失败:"+err.Error())
			fmt.Printf("读取响应体失败: %v\n", err)
			return
		}

		// 调试：打印原始 JSON (可选)
		// fmt.Printf("Raw JSON: %.200s...\n", string(body)) // 打印前200个字符

		var bubblesList []bubbles1000Data
		// 解析JSON数据
		err = json.Unmarshal(body, &bubblesList)
		if err != nil {
			global.Logger.Error(ctx, "解析JSON失败:"+err.Error())
			fmt.Printf("解析JSON失败: %v\n", err)
			fmt.Printf("响应内容片段: %.200s\n", string(body)) // 打印片段帮助调试
			return
		}

		global.Logger.Info(ctx, fmt.Sprintf("成功解析到 %d 条数据", len(bubblesList)))
		fmt.Printf("成功解析到 %d 条数据\n", len(bubblesList))

		// 如果需要清空旧数据，取消下面注释
		// if err := tx.Table("bubbles1000").Where("1=1").Delete(&bubbles1000.Bubbles1000{}).Error; err != nil {
		//	 tx.Rollback()
		//	 global.Logger.Error(ctx, "清空旧数据失败:"+err.Error())
		//	 fmt.Printf("清空旧数据失败: %v\n", err)
		//	 return
		// }

		for _, data := range bubblesList {
			// fmt.Println(data) // 调试打印

			// 将map类型转换为JSON字符串
			symbolsJson, err := json.Marshal(data.Symbols)
			if err != nil {
				global.Logger.Error(ctx, "序列化 Symbols 失败:"+err.Error())
				fmt.Printf("序列化 Symbols 失败: %v\n", err)
				// 可以选择 continue 跳过错误，或 return 终止
				continue
			}
			performanceJson, err := json.Marshal(data.Performance)
			if err != nil {
				global.Logger.Error(ctx, "序列化 Performance 失败:"+err.Error())
				fmt.Printf("序列化 Performance 失败: %v\n", err)
				// 可以选择 continue 跳过错误，或 return 终止
				continue
			}
			rankDiffsJson, err := json.Marshal(data.RankDiffs)
			if err != nil {
				global.Logger.Error(ctx, "序列化 RankDiffs 失败:"+err.Error())
				fmt.Printf("序列化 RankDiffs 失败: %v\n", err)
				// 可以选择 continue 跳过错误，或 return 终止
				continue
			}
			exchangePricesJson, err := json.Marshal(data.ExchangePrices)
			if err != nil {
				global.Logger.Error(ctx, "序列化 ExchangePrices 失败:"+err.Error())
				fmt.Printf("序列化 ExchangePrices 失败: %v\n", err)
				// 可以选择 continue 跳过错误，或 return 终止
				continue
			}

			// 检查记录是否已存在，决定是创建还是更新
			var existingRecord bubbles1000.Bubbles1000
			err = global.DB.Table("bubbles1000").Where("slug = ?", data.Slug).First(&existingRecord).Error

			bubbles1000Model := bubbles1000.Bubbles1000{
				// Id:  对于更新，GORM 通常会根据主键自动处理
				Name:      data.Name,
				Slug:      data.Slug,
				Symbol:    data.Symbol,
				Dominance: data.Dominance,
				Image:     data.Image,
				Rank:      int32(data.Rank),
				Price:     data.Price,
				Marketcap: data.Marketcap,
				Volume:    data.Volume,
				CgId:      data.CgId,
				// Stable:         data.Stable, // 如果表结构中添加了此字段
				Symbols:        string(symbolsJson),
				Performance:    string(performanceJson),
				RankDiffs:      string(rankDiffsJson),
				ExchangePrices: string(exchangePricesJson),
			}

			var result *gorm.DB
			if err == nil && existingRecord.Id > 0 {
				// 记录存在，执行更新
				// 注意：这里只更新非零值字段。如果需要强制更新空值，需要使用 Select
				result = global.DB.Table("bubbles1000").Where("id = ?", existingRecord.Id).Updates(&bubbles1000Model)
				// 或者更明确地指定要更新的字段:
				// result = tx.Table("bubbles1000").Where("id = ?", existingRecord.Id).Select("*").Updates(&bubbles1000Model)
			} else {
				// 记录不存在或查询出错（假设是不存在），执行创建
				// 注意：如果因为唯一索引冲突导致创建失败，需要特殊处理
				result = global.DB.Table("bubbles1000").Create(&bubbles1000Model)
			}

			if result.Error != nil {
				global.Logger.Error(ctx, fmt.Sprintf("操作数据失败 (Slug: %s): %v", data.Slug, result.Error))
				fmt.Printf("Error operating data (Slug: %s): %v\n", data.Slug, result.Error)
				// 可以选择 continue 跳过错误，或 return 终止
				// return
				continue // 示例：跳过单条错误
			}
			// global.Logger.Info(ctx, fmt.Sprintf("成功操作数据: %s", data.Name))
			// fmt.Printf("成功操作: %s\n", data.Name)

			// 调试用 sleep，生产环境应移除
			// if global.Env == "dev" {
			//	 time.Sleep(time.Second * 2)
			// }

		}

		global.Logger.Info(ctx, "所有数据处理并插入/更新成功")
		fmt.Println("所有数据处理并插入/更新成功")

	},
}

func init() {
	// 环境参数
	Bubbles1000Cmd.Flags().StringVarP(&global.Env, "env", "e", "dev", "测试环境")
}
