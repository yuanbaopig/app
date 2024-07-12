package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yuanbaopig/app"
	"github.com/yuanbaopig/app/fname"
)

type options struct {
	RedisOption *RedisOption `json:"Redis"`
}

func (o *options) Flags() (fss fname.NamedFlagSets) {

	o.RedisOption.AddFlags(fss.FlagSet("redis"))
	return fss
}

func (o *options) Validate() []error {
	var errs []error

	errs = append(errs, o.RedisOption.Validate()...)

	return errs
}

type RedisOption struct {
	host *string
}

func (r *RedisOption) AddFlags(fs *pflag.FlagSet) {

	fs.String("redis.host", "127.0.0.1", ""+
		"redis service host address. If left blank, the following related mysql options will be ignored.")

}

func (r *RedisOption) Validate() []error {
	var errs []error

	return errs
}

func main() {

	/*
		opts := options.NewOptions()
		application := app.NewApp("IAM API Server",		// 应用的简短描述
		    basename,		// basename 会根据给定的名称被格式化为不同操作系统下的可执行文件，在linux没有啥区别，在windows下会增加.exe文件扩展名
		    app.WithOptions(opts),
		    app.WithDescription(commandDesc),		// 应用描述
		    app.WithDefaultValidArgs(),
		    app.WithRunFunc(run(opts)),
		)
	*/

	redisFunc := func(args []string) error {
		fmt.Println("redis sub command")
		return nil
	}

	o := &options{
		RedisOption: &RedisOption{},
	}

	redisCmd := app.NewCommand("redis", "redis",
		app.WithCommandOptions(o),
		app.WithCommandRunFunc(redisFunc),
	)

	a := app.NewApp("test", "test",
		app.WithNoVersion(),
		app.WithNoConfig(),
		app.WithDescription("commandDesc"),
		app.WithAddCommand(redisCmd),
		app.WithRunFunc(func(cmd *cobra.Command, args []string) error {
			fmt.Println("root command")
			return nil
		}))

	a.Run()

}
