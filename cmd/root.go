/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type arrayFlags []string

var (
	// Used for flags.
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "log2curl",
	Short: "log2curl is a CLI that helps converting log to CURL script",
	Long:  `It's easy to use log2curl to convert the log on the cloud services (e.g., Cloud Watch) to CURL script`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.log2curl.json)")
	rootCmd.PersistentFlags().StringP("out", "o", "", "The output file <file> (default is ./log2curl.{timestamp}.txt)")
	viper.BindPFlag("output.path", rootCmd.PersistentFlags().Lookup("out"))

	rootCmd.PersistentFlags().StringP("path", "", "path", "The key represents the path of log")
	viper.BindPFlag("key.path", rootCmd.PersistentFlags().Lookup("path"))
	rootCmd.PersistentFlags().StringP("host", "", "host", "The key represents the host of log")
	viper.BindPFlag("key.host", rootCmd.PersistentFlags().Lookup("host"))
	rootCmd.PersistentFlags().StringP("method", "", "method", "The key represents the HTTP method of log")
	viper.BindPFlag("key.method", rootCmd.PersistentFlags().Lookup("method"))
	rootCmd.PersistentFlags().StringP("scheme", "", "scheme", "The key represents the HTTP scheme of log")
	viper.BindPFlag("key.scheme", rootCmd.PersistentFlags().Lookup("scheme"))
	rootCmd.PersistentFlags().StringP("headers-prefix", "", "headers", "The key prefix represents the headers of log")
	viper.BindPFlag("key.headers_prefix", rootCmd.PersistentFlags().Lookup("headers-prefix"))
	rootCmd.PersistentFlags().StringP("parameters-prefix", "", "parameters", "The key prefix represents the parameters of log")
	viper.BindPFlag("key.parameters_prefix", rootCmd.PersistentFlags().Lookup("parameters-prefix"))
	rootCmd.PersistentFlags().StringSlice("whitelist-headers", []string{"Content-Type", "Authorization"}, "The headers is used for CURL request headers")
	viper.BindPFlag("whitelist_headers", rootCmd.PersistentFlags().Lookup("whitelist-headers"))
	rootCmd.PersistentFlags().StringP("custom-host", "", "", "Custom CURL request host")
	viper.BindPFlag("custom.host", rootCmd.PersistentFlags().Lookup("custom-host"))

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".log2curl")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
