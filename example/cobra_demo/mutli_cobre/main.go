package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile     string
	projectBase string
	userLicense string
)

func main() {
	rootCmd := &cobra.Command{
		Use: "myapp",
		// Short is the short description shown in the 'help' output.
		Short: "Hugo is a very fast static site generator",
		// Long is the long message shown in the 'help <this-command>' output.
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 主命令的逻辑
			fmt.Println("Executing main command")
			fmt.Println(viper.GetString("globalFlag"))
			return nil
		},
	}

	//cobra.OnInitialize(initConfig)
	// 在父命令中设置选项
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	rootCmd.PersistentFlags().StringP("globalFlag", "g", "", "A global flag for the main command")

	// viper set
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")

	subCmd := &cobra.Command{
		Use:   "sub",
		Short: "A subcommand",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 子命令的逻辑
			fmt.Println("Executing subcommand")
			fmt.Println(viper.GetString("exampleFlag"))
			fmt.Println(viper.GetString("globalFlag"))
			return nil
		},
	}

	// 在子命令中设置选项
	subCmd.PersistentFlags().StringP("exampleFlag", "e", "", "An example flag for the subcommand")
	// 绑定整个pflag
	viper.BindPFlags(subCmd.PersistentFlags())

	rootCmd.AddCommand(subCmd)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

// initConfig read config
func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	// viper read config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
