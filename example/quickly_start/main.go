package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yuanbaopig/app"
)

func main() {
	var port string
	pflag.StringVarP(&port, "port", "p", "1", "help info")

	app.NewApp("app", "basename",
		//app.WithSilence(),
		app.WithValidArgs(cobra.MinimumNArgs(1)),
		app.WithNoConfig(),
		app.WithDescription("description"),
		app.WithNoVersion(),
		app.WithRunFunc(func(cmd *cobra.Command, args []string) error {
			fmt.Println(port)
			fmt.Println(args)
			return nil
		}),
	).Run()

}
