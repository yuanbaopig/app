# 项目名称

项目代码来源于：github.com/marmotedu/component-base/pkg/app

应用构建包，具备命令行程序、命令行参数解析和配置文件解析这 3 种功能。

- **命令行程序**：用来启动一个应用。命令行程序需要实现诸如应用描述、help、参数校验等功能。
- **命令行参数解析**：用来在启动时指定应用程序的命令行参数，以控制应用的行为。
- **配置文件解析**：用来解析不同格式的配置文件。



该包具有如下特性：

- 支持自定义分组选项参数
- 支持多层级的选项参数
- 支持多种配置文件解析
- 支持命令行无选项参数



## 功能特性

- `WithNoConfig()`：不指定配置文件，配置文件支持默认路径和指定文件
- `WithValidArgs(args cobra.PositionalArgs)`：用户命令行无选项参数
- `WithOptions(opt CliOptions) `：用户自定义分组选项参数
-  `WithDescription(desc string)`：用户命令描述
- `WithCommands(commands ...*cobra.Command)`：用户多层级的选项参数





## 快速开始

开箱即用

```go
package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/yuanbaopig/app"
)

func main() {
	var port string
	pflag.StringVarP(&port, "port", "p", "1", "help info")  // 顶级选项参数

	app.NewApp("app", "basename",
		app.WithSilence(),			// 静默模式，不输出配置文件与版本相关信息
		app.WithValidArgs(cobra.MinimumNArgs(1)),			// 命令行参数校验函数，支持自定义
             // MinimumNArgs(1) 的作用是必须在命令行后面输入一个参数，注意不是选项
		app.WithNoConfig(),			// 不涉及配置文件
		app.WithDescription("description"),		// app 应用描述
		app.WithNoVersion(),		// 不涉及build version相关信息，需要单独维护
		app.WithRunFunc(func(cmd *cobra.Command, args []string) error {   // 运行的app
			fmt.Println(port)
			fmt.Println(args)			// 命令行参数调用
			return nil
		}),
	).Run()

}
```

默认配置文件的名称为basename+扩展字段，例如basename为api-server，配置文件则为api-server.ini。支持多种类型的格式，扩展名根据对应类型添加。



## 使用指南

### 用户选项分组

```go
package main

import (
	"github.com/spf13/pflag"
	"github.com/yuanbaopig/app"
	"github.com/yuanbaopig/app/fname"
)

// 选项组
type options struct {
	MySQLOptions *MySQLOptions `json:"MySQL"`
	RedisOption  *RedisOption  `json:"Redis"`
}

// 初始化中需要调用的一个方法
func (o *options) Flags() (fss fname.NamedFlagSets) {
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOption.AddFlags(fss.FlagSet("redis"))
	return fss
}

func (o *options) Validate() []error {
	var errs []error

	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOption.Validate()...)

	return errs
}

type MySQLOptions struct {
	host *string
}

func (m *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	// 正式环境中也可以先new一个MySQLOptions，然后将返回值存储在m.host 中
	fs.String("mysql.host", "127.0.0.1", ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

}

func (m *MySQLOptions) Validate() []error {
	var errs []error
    // 此处是用来校验参数合法性的，如果返回值不是空，应用启动就会报错
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
	// new optins
	o := &options{
		MySQLOptions: &MySQLOptions{},
		RedisOption:  &RedisOption{},
	}

	app.NewApp("test", "test",
		app.WithSilence(),
		app.WithNoVersion(),
		//app.WithNoConfig(),				// 想要使用viper功能，必须开启config
		app.WithOptions(o),
		app.WithRunFunc(run),
	).Run()

}

// 启用配置文件，则可以绑定viper
func run(cmd *cobra.Command, args []string) error {
	fmt.Println("test")
	fmt.Println(viper.GetString("mysql.host"))
	fmt.Println(viper.GetString("redis.host"))		// viper 选项参数绑定调用
	fmt.Println(args)
	return nil
}
```

### 多命令选项

用户命令多层级选项参数，并且进行viper绑定

```go
package main

import (
	"fmt"
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
```

