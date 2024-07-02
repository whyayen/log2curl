package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/magiconair/properties/assert"
	"github.com/whyayen/log2curl/internal/models"
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
			"host":                 "example.com",
			"header.Authorization": "Bearer token",
			"header.Content-Type":  "application/json",
			"header.User-Agent":    "Webkit 1.0",
			"parameter.id":         "66838a1d337a8cdc830b439c",
			"parameter.subscribed": "true",
			"path":                 "/users",
			"method":               "GET",
			"scheme":               "https",
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
	svc.RestfulConfiguration = &models.RestfulConfiguration{
		Host:             "host",
		Path:             "path",
		Method:           "method",
		Scheme:           "scheme",
		HeaderPrefix:     "header",
		ParameterPrefix:  "parameter",
		WhitelistHeaders: []string{"Authorization"},
		CustomHost:       "custom.example2.com",
	}

	results, err := svc.GetQueryResults(ctx)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(results.Results), 1)
	assert.Equal(t, results.Results[0].GetHost(), "custom.example2.com")
	assert.Equal(t, results.Results[0].GetPath(), "/users")
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
