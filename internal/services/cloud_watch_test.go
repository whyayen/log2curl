package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/magiconair/properties/assert"
	"testing"
)

type mockCloudWatchLogsClient struct {
	cloudwatchlogs.Client
}

func (m *mockCloudWatchLogsClient) GetQueryResults(ctx context.Context, params *cloudwatchlogs.GetQueryResultsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	var results [][]types.ResultField
	row1 := []types.ResultField{
		types.ResultField{
			Field: func(s string) *string { return &s }("field"),
			Value: func(s string) *string { return &s }("value"),
		},
		types.ResultField{
			Field: func(s string) *string { return &s }("field2"),
			Value: func(s string) *string { return &s }("value2"),
		},
		types.ResultField{
			Field: func(s string) *string { return &s }("field3"),
			Value: func(s string) *string { return &s }("value3"),
		},
	}
	results = append(results, row1)
	return &cloudwatchlogs.GetQueryResultsOutput{
		Status:  "Complete",
		Results: results,
	}, nil
}

func (m *mockCloudWatchLogsClient) GetLogRecord(ctx context.Context, params *cloudwatchlogs.GetLogRecordInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogRecordOutput, error) {
	return &cloudwatchlogs.GetLogRecordOutput{
		LogRecord: map[string]string{
			"key":  "value",
			"key2": "value2",
		},
	}, nil
}

type mockCloudWatchReturnErrClient struct {
	cloudwatchlogs.Client
}

func (m *mockCloudWatchReturnErrClient) GetQueryResults(ctx context.Context, params *cloudwatchlogs.GetQueryResultsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	return nil, &types.InvalidParameterException{}
}

func TestGetQueryResults(t *testing.T) {
	ctx := context.TODO()
	queryId := "queryId"
	svc := CloudWatchService{
		CloudWatchQueryParams: CloudWatchQueryParams{
			queryId: &queryId,
		},
		client: &mockCloudWatchLogsClient{},
	}

	results, err := svc.GetQueryResults(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(results.Results), 1)
}

func TestGetQueryResultsWithError(t *testing.T) {
	ctx := context.TODO()
	queryId := "queryId"
	svc := CloudWatchService{
		CloudWatchQueryParams: CloudWatchQueryParams{
			queryId: &queryId,
		},
		client: &mockCloudWatchReturnErrClient{},
	}

	results, err := svc.GetQueryResults(ctx)

	assert.Equal(t, err == nil, false)
	assert.Equal(t, len(results.Results), 0)
}