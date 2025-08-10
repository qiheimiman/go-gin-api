package command

import (
	"github.com/xinliangnote/go-gin-api/cmd/command/binance"
	"github.com/xinliangnote/go-gin-api/cmd/command/bubbles"
	"github.com/xinliangnote/go-gin-api/cmd/command/zhonghuan"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(zhonghuan.ZhonghuanCmd)
	rootCmd.AddCommand(bubbles.Bubbles1000Cmd)
	rootCmd.AddCommand(binance.BinanceNewsCmd)
	rootCmd.AddCommand(binance.BinanceFearGreedCmd)

}
