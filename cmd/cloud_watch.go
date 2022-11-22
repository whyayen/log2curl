package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/whyayen/log2curl/services"
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
		fmt.Println("Start to get query results...")
		services.CloudWatch(&queryId)
	},
}
