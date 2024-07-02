package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config bool

func init() {
	generateCmd.Flags().BoolVarP(&config, "config", "", false, "generate config file to $HOME/.log2curl.json")
	generateCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate some file",
	Run: func(cmd *cobra.Command, args []string) {
		if config {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			viper.AddConfigPath(home)
			viper.SetConfigName(".log2curl")
			viper.SetConfigType("json")
			viper.SetDefault("custom.host", "")
			viper.SetDefault("key.path", "path")
			viper.SetDefault("key.host", "host")
			viper.SetDefault("key.method", "method")
			viper.SetDefault("key.scheme", "scheme")
			viper.SetDefault("key.headers_prefix", "headers")
			viper.SetDefault("key.parameters_prefix", "parameters")
			viper.SetDefault("whitelist_headers", []string{"Content-Type", "Authorization"})
			err = viper.SafeWriteConfig()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Generate default config to $HOME/.log2curl.json successfully")
		}
	},
}
