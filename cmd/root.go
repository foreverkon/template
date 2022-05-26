package cmd

import (
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// some color functions
var (
	blue    = color.New(color.FgBlue).SprintFunc()
	green   = color.New(color.FgGreen)
	red     = color.New(color.FgRed)
	yellow  = color.New(color.FgYellow)
	magenta = color.New(color.FgMagenta)
)

var rootCmd = &cobra.Command{
	Use:   "template <cmd> <flag> <option>",
	Short: "A command-line  tool to download some template files",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// define global flags
	rootCmd.PersistentFlags().StringP("name", "n", "", "name of the template to search or download")

	// define local flags
	rootCmd.Flags().StringP("templates", "t", "", "path of the tempaltes.json")

	// bind flags to viper
	viper.BindPFlag("templates", rootCmd.LocalFlags().Lookup("templates"))
}

func initConfig() {
	if viper.GetString("templates") != "" {
		viper.SetConfigFile(viper.GetString("templates"))
	} else {
		viper.AutomaticEnv()
		viper.SetConfigFile(viper.GetString("GOPATH") + "\\cfg\\templates.json")
	}
	err := viper.ReadInConfig()
	if err != nil {
		// if there is no config file, create one
		viper.SetDefault("update", time.Now().Unix())
		viper.SetDefault("names", []string{})
		viper.SetDefault("tempaltes", map[string]string{})
		viper.WriteConfig()
	}
}
