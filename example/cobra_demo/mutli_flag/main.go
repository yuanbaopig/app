package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	// 创建两个 FlagSet 对象
	fs1 := pflag.NewFlagSet("fs1", pflag.ExitOnError)
	fs2 := pflag.NewFlagSet("fs2", pflag.ExitOnError)

	// 向第一个 FlagSet 添加标记
	var flag1 string
	fs1.StringVar(&flag1, "flag1", "default1", "Description for flag1")

	// 将第一个 FlagSet 添加到第二个 FlagSet 中
	fs2.AddFlagSet(fs1)

	// 在第二个 FlagSet 中添加额外的标记
	var flag2 int
	fs2.IntVar(&flag2, "flag2", 42, "Description for flag2")

	// 创建一个命令
	cmd := &cobra.Command{
		Use:   "myapp",
		Short: "A simple example app",
		Run: func(cmd *cobra.Command, args []string) {
			// 输出标记值
			fmt.Println("Flag1:", flag1)
			fmt.Println("Flag2:", flag2)
		},
	}

	// 将第二个 FlagSet 添加到命令的 FlagSet 中
	cmd.Flags().AddFlagSet(fs2)

	// 执行命令
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
