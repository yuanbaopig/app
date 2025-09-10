// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/marmotedu/errors"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yuanbaopig/app/fname"
	"github.com/yuanbaopig/app/version"
	"github.com/yuanbaopig/app/version/verflag"
	"io"
	"os"
	"strings"
)

var (
	progressMessage = color.GreenString("==>")

	usageTemplate = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

// App is the main structure of a cli application.
// It is recommended that an app be created with the app.NewApp() function.
type App struct {
	basename    string
	name        string
	description string
	options     CliOptions
	runFunc     RunFunc
	silence     bool
	noVersion   bool
	noConfig    bool
	commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

// Option defines optional parameters for initializing the application
// structure.
type Option func(*App)

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithSilence sets the application to silent mode, in which the program startup
// information, configuration information, and version information are not
// printed in the console.
func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoVersion set the application does not provide version flag.
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// WithNoConfig set the application does not provide config flag.
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithValidArgs set the validation function to valid non-flag arguments.
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

// NewApp creates a new application instance based on the given application name,
// binary name, and other options.
func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:   FormatBaseName(a.basename),
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	// cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	// 修改flags 选项名称中的符号
	fname.InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}
	// 运行函数赋值
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	// 将flag Set中的pflag添加到cmd flags中
	var namedFlagSets fname.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		//fs := cmd.Flags()
		fs := cmd.PersistentFlags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}
	// 检查是否设置了version选项，如果设置了则添加对应的选项参数，默认设置
	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}
	// 检查是否设置了config选项，如果设置了则添加对应的选项参数，默认设置
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}
	// 配置help选项信息
	AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	// add new global flagset to cmd FlagSet
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))
	// 设置自定义使用信息和帮助信息
	if a.commands == nil {
		addCmdTemplate(&cmd, namedFlagSets)
	}

	a.cmd = &cmd
}

// Run is used to launch the application.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// RunContext is used to launch the application with context.
func (a *App) RunContext(ctx context.Context) {
	if err := a.cmd.ExecuteContext(ctx); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command returns cobra command instance inside the application.
func (a *App) Command() *cobra.Command {
	return a.cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	var pb strings.Builder // 记录viper映射config前的flags建值
	//var afterConfig []string
	if !a.silence { // config配置必须在bind flags之前打印
		if !a.noConfig {
			fmt.Printf("%v Config file used: `%s`\n", progressMessage, viper.ConfigFileUsed())
			afterConfig := viper.AllKeys() // 配置文件的建值
			printConfig(afterConfig)
		}

	}

	if !a.noConfig {
		//afterConfig = viper.AllKeys() // 配置文件的建值

		pbF := func(flag *pflag.Flag) {
			//pb.WriteString(fmt.Sprintf("FLAG: --%s=%q\n", flag.Name, flag.Value))
			if flag.Changed || flag.Value.String() != flag.DefValue {
				// 只有在 flag 被设置，或者值与默认值不同的情况下才调用 pb.WriteString
				pb.WriteString(fmt.Sprintf("FLAG: --%s=%q\n", flag.Name, flag.Value))
			}
		}
		cmd.Flags().VisitAll(pbF)

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		printWorkingDir()
		fmt.Printf("%v Flags items:\n", progressMessage)
		fmt.Printf(pb.String())

		fmt.Printf("%v Starting %s ...\n", progressMessage, a.name)
		if !a.noVersion {
			fmt.Printf("%v Version: `%s`\n", progressMessage, version.Get().ToJSON())
		}
		//if !a.noConfig {
		//	fmt.Printf("%v Config file used: `%s`\n", progressMessage, viper.ConfigFileUsed())
		//	printConfig(afterConfig)
		//}
	}
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	// 首先检查 a.options 是否实现了 CompletableOptions 接口。
	// Go语言中，接口的实现是隐式的，我们可以通过类型断言来判断某个变量是否实现了某个接口。
	if CompletableOption, ok := a.options.(CompletableOptions); ok {
		// 如果 a.options 实现了 CompletableOptions 接口，那么就调用这个接口的 Complete 方法。
		// 完成之后，检查是否有错误发生，如果有错误，那么就直接返回这个错误。
		if err := CompletableOption.Complete(); err != nil {
			return err
		}
	}
	// 调用 a.options 的 Validate 方法，该方法返回一个包含所有错误的切片。
	// 如果返回的错误切片的长度不为0，表明验证过程中出现了错误，那么创建一个新的错误聚合并返回。
	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}
	// 检查 a.options 是否实现了 PrintableOptions 接口，并且 App 是否设置为禁声模式（a.silence 不为 true）。
	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		// 如果实现了 PrintableOptions 接口且 App 没有被设置为禁声模式，
		// 那么就打印 options 的配置信息。这里假设 progressMessage 是个已定义的全局变量。
		fmt.Printf("%v Config: `%s`\n", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	fmt.Printf("%v WorkingDir: %s\n", progressMessage, wd)
}

// 自定义格式的使用和帮助信息
func addCmdTemplate(cmd *cobra.Command, namedFlagSets fname.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := TerminalSize(cmd.OutOrStdout())
	// 设置命令的使用信息打印函数
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		fname.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	// 设置命令的帮助信息打印函数
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		fname.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}

func TerminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winSize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winSize.Width), int(winSize.Height), nil
}
