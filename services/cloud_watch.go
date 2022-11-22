package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/viper"
	"github.com/whyayen/log2curl/convertors"
	"github.com/whyayen/log2curl/parsers"
)

func CloudWatch(queryId *string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := cloudwatchlogs.NewFromConfig(cfg)
	output, err := client.GetQueryResults(context.TODO(), &cloudwatchlogs.GetQueryResultsInput{
		QueryId: queryId,
	})

	if output.Status != "Complete" {
		fmt.Println("The query isn't complete. Please run again after the query completing.")
		os.Exit(1)
	}

	// collect @ptr
	var ptrSlice []*string
	for _, result := range output.Results {
		ptrSlice = append(ptrSlice, result[2].Value)
	}

	config := parsers.ConfigMapperInterface{
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

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	total := len(ptrSlice)
	fmt.Printf("Processing: %d / %d \r", 0, total)
	for index, ptr := range ptrSlice {
		response, err := client.GetLogRecord(context.TODO(), &cloudwatchlogs.GetLogRecordInput{
			LogRecordPointer: ptr,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		payload := parsers.General(&config, &response.LogRecord)

		if payload != "" {
			f.WriteString(convertors.Curl(&payload) + "\n\n")
		} else {
			f.WriteString("Failed\n\n")
		}

		fmt.Printf("Processing: %d / %d \r", index+1, total)
	}

	fmt.Printf("Finished. Save file in: %s", filePath)
}
