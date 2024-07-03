package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/viper"
	"github.com/whyayen/log2curl/internal/models"
	"github.com/whyayen/log2curl/internal/parsers"
	"sync"
)

type CloudWatchQueryParams struct {
	queryId *string
}

type CloudWatchService struct {
	CloudWatchQueryParams
	*models.HttpRequestConfiguration
	client CloudWatchClient
}

type CloudWatchResults struct {
	Results []*models.HttpRequest
}

type CloudWatchClient interface {
	GetQueryResults(ctx context.Context, params *cloudwatchlogs.GetQueryResultsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetQueryResultsOutput, error)
	GetLogRecord(ctx context.Context, params *cloudwatchlogs.GetLogRecordInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogRecordOutput, error)
}

func NewCloudWatchService(ctx context.Context, queryId *string) (*CloudWatchService, error) {
	awsProfile := viper.GetString("aws.profile")
	var cfg aws.Config
	var err error

	if awsProfile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(awsProfile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return &CloudWatchService{}, fmt.Errorf("error loading default config: %w", err)
	}

	client := cloudwatchlogs.NewFromConfig(cfg)

	return &CloudWatchService{
		CloudWatchQueryParams: CloudWatchQueryParams{
			queryId: queryId,
		},
		client: client,
	}, nil
}

func (c *CloudWatchService) GetQueryResults(ctx context.Context) (CloudWatchResults, error) {
	output, err := c.client.GetQueryResults(ctx, &cloudwatchlogs.GetQueryResultsInput{
		QueryId: c.queryId,
	})

	if err != nil {
		return CloudWatchResults{}, fmt.Errorf("error getting query results: %w", err)
	}

	if output.Status != "Complete" {
		return CloudWatchResults{}, fmt.Errorf("query isn't complete. Please run again after the query completing")
	}

	var ptrSlice []*string
	for _, result := range output.Results {
		for index, field := range result {
			if *field.Field == "@ptr" {
				ptrSlice = append(ptrSlice, result[index].Value)
			}
		}
	}

	var wg sync.WaitGroup
	resultsChan := make(chan *models.HttpRequest, len(ptrSlice))
	errChan := make(chan error, len(ptrSlice))
	parser := parsers.NewCloudWatchParser(c.HttpRequestConfiguration)
	wg.Add(len(ptrSlice))

	for _, ptr := range ptrSlice {
		go func(ptr *string) {
			defer wg.Done()
			resp, err := c.client.GetLogRecord(ctx, &cloudwatchlogs.GetLogRecordInput{
				LogRecordPointer: ptr,
			})

			if err != nil {
				errChan <- err
				return
			}

			resultsChan <- parser.Parse(&resp.LogRecord)
		}(ptr)
	}

	wg.Wait()
	close(resultsChan)
	close(errChan)

	if len(errChan) > 0 {
		return CloudWatchResults{}, fmt.Errorf("error getting log record: %w", <-errChan)
	}

	results := CloudWatchResults{
		Results: []*models.HttpRequest{},
	}
	for result := range resultsChan {
		results.Results = append(results.Results, result)
	}

	return results, nil
}
