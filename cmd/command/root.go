package command

import (
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

}
