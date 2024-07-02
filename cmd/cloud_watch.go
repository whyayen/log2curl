package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/whyayen/log2curl/internal/generators"
	"github.com/whyayen/log2curl/internal/models"
	"github.com/whyayen/log2curl/internal/services"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var queryId string

func init() {
	cloudWatchCmd.Flags().StringVarP(&queryId, "query-id", "q", "", "Cloud Watch query ID (required)")
	cloudWatchCmd.MarkFlagRequired("query-id")

	rootCmd.AddCommand(cloudWatchCmd)
}

var cloudWatchCmd = &cobra.Command{
	Use:     "cloud_watch",
	Aliases: []string{"cw"},
	Short:   "Convert the log on Cloud Watch to CURL",
	Long: `It's easy to convert the log on Cloud Watch to CURL command with Query Id.
    The query id must be given when running this command. And the status of the query should be Complete.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Loading default config...")
		requestCfg := models.HttpRequestConfiguration{
			ParameterPrefix:  viper.GetString("key.parameters_prefix"),
			HeaderPrefix:     viper.GetString("key.headers_prefix"),
			Path:             viper.GetString("key.path"),
			Host:             viper.GetString("key.host"),
			Method:           viper.GetString("key.method"),
			Scheme:           viper.GetString("key.scheme"),
			WhitelistHeaders: viper.GetStringSlice("whitelist_headers"),
			CustomHost:       viper.GetString("custom.host"),
		}
		var filePath string
		if viper.GetString("output.path") != "" {
			filePath = viper.GetString("output.path")
		} else {
			filePath = fmt.Sprintf("log2curl.%d.txt", time.Now().Unix())
		}

		fmt.Println("Start to get query results...")
		ctx := context.TODO()
		svc, err := services.NewCloudWatchService(ctx, &queryId)
		svc.HttpRequestConfiguration = &requestCfg

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		results, err := svc.GetQueryResults(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Get query results successfully!")

		fmt.Println("Start to generate CURL command...")
		f, err := os.Create(filePath)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		total := len(results.Results)
		fmt.Printf("Processing: %d / %d \r", 0, total)

		for index, result := range results.Results {
			command, err := generators.NewCurlGenerator(result).Generate()

			if err != nil {
				command = err.Error()
			}

			_, err = f.WriteString(command + "\n\n")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Processing: %d / %d \r", index+1, total)
		}

		fmt.Printf("Finished. Save file in: %s", filePath)
	},
}
