package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yuanbaopig/app"
	"github.com/yuanbaopig/app/fname"
)

type options struct {
	MySQLOptions *MySQLOptions `json:"MySQL"`
	RedisOption  *RedisOption  `json:"Redis"`
}

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

func (o *options) Complete() error {
	return o.RedisOption.Complete()
}

type MySQLOptions struct {
	host *string
}

func (m *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String("mysql.host", "127.0.0.1", ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

}

func (m *MySQLOptions) Validate() []error {
	var errs []error

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

func (r *RedisOption) Complete() error {
	fmt.Println("call complete")
	return fmt.Errorf("complete error")
}

func main() {
	//namedFlagSets := &app.NamedFlagSets{}
	//nfs := namedFlagSets.FlagSet("global")
	//app.AddGlobalFlags(nfs, "test-cmd")
	// 创建两个 FlagSet 对象
	o := &options{
		MySQLOptions: &MySQLOptions{},
		RedisOption:  &RedisOption{},
	}

	app.NewApp("test", "config",
		//app.WithSilence(),
		app.WithNoVersion(),
		app.WithNoConfig(),
		app.WithOptions(o),
		app.WithRunFunc(run),
	).Run()

}

func run(cmd *cobra.Command, args []string) error {
	fmt.Println("test")
	fmt.Println(viper.GetString("mysql.host"))
	fmt.Println(viper.GetString("redis.host"))
	fmt.Println(args)
	return nil
}
